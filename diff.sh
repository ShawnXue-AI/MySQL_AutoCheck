#!/usr/bin/env bash

usage() {
    cat << EOF
用法: diff.sh <目录A> <目录B>

比较两个目录中每个 IP_端口 的文件数量。
输出目录A中缺失、以及数量多于目录B的 IP_端口 及其文件列表。

文件名示例：
  mysql_inspection_192.168.122.1_3310_20260421_162250.log
  MySQL_inspection_172.24.7.94_3306_20260424.docx
EOF
    exit 1
}

if [ $# -ne 2 ]; then
    usage
fi

DIR_A="$1"
DIR_B="$2"

for dir in "$DIR_A" "$DIR_B"; do
    if [ ! -d "$dir" ]; then
        echo "错误: 目录 '$dir' 不存在或不是目录" >&2
        exit 1
    fi
done

# 提取文件名中的 "IP_端口" (如 192.168.1.1_3306)
get_key() {
    local filename
    filename=$(basename "$1")
    echo "$filename" | grep -oE '[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}_[0-9]+' | head -1
}

# 声明关联数组：存储每个 key 对应的文件列表（用于最后展示多余的文件）
declare -A files_in_A   # key => 文件列表（用 | 分隔）
declare -A files_in_B
declare -A count_A
declare -A count_B

# 扫描目录A
for f in "$DIR_A"/*; do
    [ -f "$f" ] || continue
    key=$(get_key "$f")
    if [ -n "$key" ]; then
        ((count_A["$key"]++))
        files_in_A["$key"]="${files_in_A["$key"]}|$(basename "$f")"
    else
        echo "警告: 文件 '$f' 中未找到 IP_端口 模式，已忽略" >&2
    fi
done

# 扫描目录B
for f in "$DIR_B"/*; do
    [ -f "$f" ] || continue
    key=$(get_key "$f")
    if [ -n "$key" ]; then
        ((count_B["$key"]++))
        files_in_B["$key"]="${files_in_B["$key"]}|$(basename "$f")"
    else
        echo "警告: 文件 '$f' 中未找到 IP_端口 模式，已忽略" >&2
    fi
done

# 收集所有 key（A 和 B 的并集）
all_keys=()
for key in "${!count_A[@]}"; do
    all_keys+=("$key")
done
for key in "${!count_B[@]}"; do
    # 避免重复添加
    found=0
    for k in "${all_keys[@]}"; do
        [ "$k" = "$key" ] && found=1 && break
    done
    [ $found -eq 0 ] && all_keys+=("$key")
done

# 输出结果
missing_any=false
for key in "${all_keys[@]}"; do
    cntA=${count_A["$key"]:-0}
    cntB=${count_B["$key"]:-0}

    if [ $cntA -gt 0 ] && [ $cntB -eq 0 ]; then
        echo "【缺失】IP_端口 $key 存在于目录A，但目录B中完全没有对应文件。"
        echo "  目录A中的文件："
        IFS='|' read -ra files <<< "${files_in_A["$key"]}"
        for file in "${files_a[@]}"; do
            [ -n "$file" ] && echo "    - $file"
        done
        missing_any=true
    elif [ $cntA -gt $cntB ]; then
        echo "【数量不足】IP_端口 $key：目录A有 $cntA 个文件，目录B只有 $cntB 个。"
        echo "  目录A中多出的文件（前 $((cntA - cntB)) 个视为多余，具体顺序无意义）："
        IFS='|' read -ra filesA <<< "${files_in_A["$key"]}"
        IFS='|' read -ra filesB <<< "${files_in_B["$key"]}"
        # 此处简单输出目录A中所有的文件（因为无法确定B具体匹配哪个）
        # 更精确的做法是只列出文件名不同的，但考虑到命名可能仅时间戳不同，全部列出供用户判断。
        for file in "${filesA[@]}"; do
            [ -n "$file" ] && echo "    - $file"
        done
        missing_any=true
    fi
done

if [ "$missing_any" = false ]; then
    echo "完美匹配：目录A中每个 IP_端口 的文件数量与目录B完全一致。"
fi

# 可选：输出目录B中多余的文件（如果用户关注对称性）
echo ""
echo "--- 补充：目录B中独有的 IP_端口（目录A中没有的） ---"
extra_b=false
for key in "${!count_B[@]}"; do
    if [ ${count_A["$key"]:-0} -eq 0 ]; then
        echo "IP_端口 $key 仅存在于目录B："
        IFS='|' read -ra files <<< "${files_in_B["$key"]}"
        for file in "${files[@]}"; do
            [ -n "$file" ] && echo "    - $file"
        done
        extra_b=true
    fi
done
[ "$extra_b" = false ] && echo "（无）"
