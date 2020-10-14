#proto_path=module/account/proto
#go_out=module/account/proto
#micro_out=module/account/proto
#file_name=module/account/proto/user.proto

echo proto_path
read proto_path
echo ${proto_path}

echo file_name
read file_name
echo ${file_name}
protoc --proto_path=${proto_path}   --go_out=${proto_path} --micro_out=${proto_path} ${proto_path}/${file_name}
