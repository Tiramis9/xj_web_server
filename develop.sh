
##

#scp -P 22  -r bin keys config  root@47.56.172.167:/usr/local/server/api

PROJECT_PATH="/usr/local/server/web_server/bin"
PROJECT_NAME="xj_web_server"
USER_NAME="root"
HOSTS=("47.113.94.16")
PASSWORD="YC2JeVyZXWeXu3sT"

echo "Please Input the server password: "
#read -s PASSWORD

echo '------------------build------------------'
make web-server
cp ./bin/xj_web_server ./bin/xj_web_server_new

echo '-----------------upload-----------------'
for host in ${HOSTS[@]}
do
echo ${host}
./upload.expect ./bin/${PROJECT_NAME}_new ${USER_NAME} ${host} ${PASSWORD} ${PROJECT_PATH}
if [[ "$?" != 0 ]]; then
   exit 2
fi
done

echo '------------------restart-------------------'
for host in ${HOSTS[@]}
do
echo ${host}
./restart.expect ${PROJECT_NAME} ${USER_NAME} ${host} ${PASSWORD} ${PROJECT_PATH}
done

rm -rf ./bin/xj_web_server_new