cd ${BASH_SOURCE%/*} 2>/dev/null
[ -z "$CTRL_DIR" ] && export CTRL_DIR=$(pwd)

source ../install.config
source ../functions
source ../base.rc

stand_alone () {

    check_host $target_ip || err "Host unreachable"
    sync_file $module $target_ip
    rcmd root@$target_ip "
        source $TARGET_DIR/scripts/deploy_${module}.sh;
        install_${module} $target_ip" || fail "Abort"
}

deploy () {   

    for target_ip in ${target_ips[@]};do
        check_host $target_ip || err "Host unreachable"
        sync_file $module $target_ip
        rcmd root@$target_ip "
            source $TARGET_DIR/scripts/deploy_${module}.sh;
            install_${module} ${target_ips[@]}" || fail "Abort"
    done
}

output_configuration () {
  
    if [ "$mode" == "single" ];then
        apollo_host=$target_ip
        sed -i "s#^apollo.meta1.host=.*#apollo.meta1.host=$apollo_host#;
                s#^apollo.meta1.port=.*#apollo.meta1.port=8080#;
                s#^apollo.portal.url=.*#apollo.portal.url=http://$apollo_host:8070#;
                " $CTRL_DIR/base_component.properties
    elif [ "$mode" == "cluster" ];then
        sed -i "s#^apollo.meta1.host=.*#apollo.meta1.host=${target_ips[0]}#;
                s#^apollo.meta1.port=.*#apollo.meta1.port=8080#;
                s#^apollo.meta2.host=.*#apollo.meta2.host=${target_ips[1]}#;
                s#^apollo.meta2.port=.*#apollo.meta2.port=8080#;
                s#^apollo.portal.url=.*#apollo.portal.url=http://${target_ips[0]}:8070#;
                " $CTRL_DIR/base_component.properties
    fi

}

install_apollo () {    
    local module=apollo
    local target_ips=($@)
    local db_address=$1

    log "$module: installation docker"
    if ! which docker >/dev/null 2>&1; then
        wget https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo -O /etc/yum.repos.d/docker-ce.repo || fail "installation docker failed"
        _yum docker-ce
    fi
    systemctl enable docker
    systemctl start docker

    if [ "${apollo_env}" == "pro,uat" ];then
        [ $# -eq 3 ] || fail "Abort：need 3 devices"

        if [ $lan_ip != $3 ];then
            local apollo_db=apolloconfig_prod_db
            local apollo_user=apollo_prod_user
            local apollo_pwd=$apollo_prod_pwd
            local pro_meta="http://$1:8080,http://$2:8080"
            local uat_meta="http://$3:8080"
            local pro_eureka="http://$1:8080/eureka/,http://$2:8080/eureka/"
            local uat_eureka="http://$3:8080/eureka/"

            if [ $lan_ip == "$1" ];then
                log "$module: installation mysql"
                install_db

                log "$module: init database user"
                cd /usr/local/mysql
                ./bin/mysql -uroot --socket /data/mysql/mysql.sock -e "
                ALTER USER user() IDENTIFIED BY '$apollo_root_pwd';"

                ./bin/mysql -uroot -p$apollo_root_pwd --socket /data/mysql/mysql.sock < $TARGET_DIR/apollo/sql/apolloconfig_prod_db.sql
                ./bin/mysql -uroot -p$apollo_root_pwd --socket /data/mysql/mysql.sock < $TARGET_DIR/apollo/sql/apolloconfig_a_db.sql
                ./bin/mysql -uroot -p$apollo_root_pwd --socket /data/mysql/mysql.sock < $TARGET_DIR/apollo/sql/apolloportal_prod_db.sql
                ./bin/mysql -uroot -p$apollo_root_pwd --socket /data/mysql/mysql.sock -e "
                    GRANT ALL ON apolloconfig_prod_db.* to '${apollo_user}'@'%' IDENTIFIED BY '$apollo_prod_pwd';
                    GRANT ALL ON apolloportal_prod_db.* to '${apollo_user}'@'%' IDENTIFIED BY '$apollo_prod_pwd';
                    GRANT ALL ON apolloconfig_a_db.* to 'apollo_a_user'@'%' IDENTIFIED BY '$apollo_a_pwd';
                    UPDATE apolloconfig_prod_db.ServerConfig SET Value='${pro_eureka}' WHERE Id=1;
                    UPDATE apolloconfig_a_db.ServerConfig SET Value='${uat_eureka}' WHERE Id=1;"
            fi

            install_configservice    
            install_adminservice
            install_portal

        elif [ $lan_ip == $3 ];then
            local apollo_db=apolloconfig_a_db
            local apollo_user=apollo_a_user
            local apollo_pwd=$apollo_a_pwd

            install_configservice    
            install_adminservice           
        fi
    elif [ "${apollo_env}" == "pro" ];then
        local apollo_db=apolloconfig_prod_db
        local apollo_user=apollo_prod_user
        local apollo_pwd=$apollo_prod_pwd

        if [ $# -eq 1 ];then
            local pro_meta="http://$1:8080"
            local pro_eureka="http://$1:8080/eureka/"
        elif [ $# -eq 2 ];then
            local pro_meta="http://$1:8080,http://$2:8080"
            local pro_eureka="http://$1:8080/eureka/,http://$2:8080/eureka/"        
        fi

        if [ $lan_ip == "$1" ];then
            log "$module: installation mysql"
            install_db

            log "$module: init database user"
            cd /usr/local/mysql
            ./bin/mysql -uroot --socket /data/mysql/mysql.sock -e "
            ALTER USER user() IDENTIFIED BY '$apollo_root_pwd';"

            ./bin/mysql -uroot -p$apollo_root_pwd --socket /data/mysql/mysql.sock < $TARGET_DIR/apollo/sql/apolloconfig_prod_db.sql
            ./bin/mysql -uroot -p$apollo_root_pwd --socket /data/mysql/mysql.sock < $TARGET_DIR/apollo/sql/apolloportal_prod_db.sql
            ./bin/mysql -uroot -p$apollo_root_pwd --socket /data/mysql/mysql.sock -e "
                GRANT ALL ON apolloconfig_prod_db.* to '${apollo_user}'@'%' IDENTIFIED BY '$apollo_prod_pwd';
                GRANT ALL ON apolloportal_prod_db.* to '${apollo_user}'@'%' IDENTIFIED BY '$apollo_prod_pwd';
                UPDATE apolloconfig_prod_db.ServerConfig SET Value='${pro_eureka}' WHERE Id=1;"
        fi

        install_configservice    
        install_adminservice
        install_portal           
       
    elif [ "${apollo_env}" == "fat" ];then
        local apollo_db=apolloconfig_test_db
        local apollo_user=apollo_test_user
        local apollo_pwd=$apollo_test_pwd
        local fat_meta="http://$1:8080"
        local fat_eureka="http://$1:8080/eureka/"

        log "$module: installation mysql"
        install_db

        log "$module: init database"
        cd /usr/local/mysql
        ./bin/mysql -uroot --socket /data/mysql/mysql.sock -e "
            ALTER USER user() IDENTIFIED BY '$apollo_root_pwd';"
        ./bin/mysql -uroot -p$apollo_root_pwd --socket /data/mysql/mysql.sock < $TARGET_DIR/apollo/sql/apolloconfig_test_db.sql
        ./bin/mysql -uroot -p$apollo_root_pwd --socket /data/mysql/mysql.sock < $TARGET_DIR/apollo/sql/apolloportal_prod_db.sql
        ./bin/mysql -uroot -p$apollo_root_pwd --socket /data/mysql/mysql.sock -e "
            GRANT ALL ON apolloconfig_test_db.* to '${apollo_user}'@'%' IDENTIFIED BY '$apollo_test_pwd';
            GRANT ALL ON apolloportal_prod_db.* to '${apollo_user}'@'%' IDENTIFIED BY '$apollo_test_pwd';
            UPDATE apolloconfig_test_db.ServerConfig SET Value='${fat_eureka}' WHERE Id=1;"
        
        install_configservice    
        install_adminservice
        install_portal
    fi
   
}

install_db () {

    add_user mysql
    init_dir mysql

    _yum libaio-devel numactl-libs rsync || fail "install requirments failed"

    log "$module: Unzip the mysql installation package"
    tar -zxf ${TARGET_DIR}/mysql/mysql-5.7.28.tar.gz -C /usr/local/mysql --strip-components=1 || err "Unzip failed"
    rsync -a --delete ${TARGET_DIR}/mysql/my.cnf /etc/my.cnf
    rsync -a  /usr/local/mysql/support-files/mysql.server /usr/local/mysql/bin/mysql.sh

    sed -i 's#^basedir=.*#basedir=/usr/local/mysql#;
            s#^datadir=.*#datadir=/data/mysql/#;
            s#^mysqld_pid_file_path=.*#mysqld_pid_file_path=/data/mysql/mysql.pid#
            ' /usr/local/mysql/bin/mysql.sh

    log "$module: init mysql"
    cd /usr/local/mysql

    ./bin/mysqld --initialize-insecure --user=mysql --basedir=/usr/local/mysql --datadir=/data/mysql
 
    sleep 10
    /usr/local/mysql/bin/mysql.sh start
    sleep 5
    return $?
}

install_configservice () {
    local container_name=apollo-configservice

    log "$module: installation apollo-configservice"
    docker load -i $TARGET_DIR/apollo/apollo-configservice-${apollo_version}.tar

    docker run -p 8080:8080 \
        -e SPRING_DATASOURCE_URL="jdbc:mysql://${db_address}:3306/${apollo_db}?characterEncoding=utf8" \
        -e SPRING_DATASOURCE_USERNAME=${apollo_user} -e SPRING_DATASOURCE_PASSWORD=${apollo_pwd} \
        -e EUREKA_INSTANCE_IP_ADDRESS=${lan_ip} \
        -d -v /opt/logs:/opt/logs --name $container_name apollo/apollo-configservice:${apollo_version}

}

install_adminservice () {
    local container_name=apollo-adminservice

    log "$module: installation apollo-adminservice"
    docker load -i $TARGET_DIR/apollo/apollo-adminservice-${apollo_version}.tar

    docker run -p 8090:8090 \
        -e SPRING_DATASOURCE_URL="jdbc:mysql://${db_address}:3306/${apollo_db}?characterEncoding=utf8" \
        -e SPRING_DATASOURCE_USERNAME=${apollo_user}  -e SPRING_DATASOURCE_PASSWORD=${apollo_pwd} \
        -e EUREKA_INSTANCE_IP_ADDRESS=${lan_ip} \
        -d -v /opt/logs:/opt/logs --name $container_name apollo/apollo-adminservice:${apollo_version}

}

install_portal () {
    local container_name=apollo-portal
    
    log "$module: installation apollo-portal"
    docker load -i $TARGET_DIR/apollo/apollo-portal-${apollo_version}.tar

    if [ "${apollo_env}" == "pro,uat" ];then
        docker run -p 8070:8070 \
            -e SPRING_DATASOURCE_URL="jdbc:mysql://${db_address}:3306/apolloportal_prod_db?characterEncoding=utf8" \
            -e SPRING_DATASOURCE_USERNAME=${apollo_user}  -e SPRING_DATASOURCE_PASSWORD=${apollo_pwd}  \
            -e APOLLO_PORTAL_ENVS=${apollo_env} \
            -e UAT_META=${uat_meta} -e PRO_META=${pro_meta} \
            -d -v /opt/logs:/opt/logs --name apollo-portal apollo/apollo-portal:${apollo_version}

    elif [ "${apollo_env}" == "pro" ];then
        docker run -p 8070:8070 \
            -e SPRING_DATASOURCE_URL="jdbc:mysql://${db_address}:3306/apolloportal_prod_db?characterEncoding=utf8" \
            -e SPRING_DATASOURCE_USERNAME=${apollo_user}  -e SPRING_DATASOURCE_PASSWORD=${apollo_pwd}  \
            -e APOLLO_PORTAL_ENVS=${apollo_env} \
            -e PRO_META=${pro_meta} \
            -d -v /opt/logs:/opt/logs --name apollo-portal apollo/apollo-portal:${apollo_version}

    elif [ "${apollo_env}" == "fat" ];then
        docker run -p 8070:8070 \
            -e SPRING_DATASOURCE_URL="jdbc:mysql://${db_address}:3306/apolloportal_prod_db?characterEncoding=utf8" \
            -e SPRING_DATASOURCE_USERNAME=${apollo_user}  -e SPRING_DATASOURCE_PASSWORD=${apollo_pwd}  \
            -e APOLLO_PORTAL_ENVS=${apollo_env} \
            -e FAT_META=${fat_meta} \
            -d -v /opt/logs:/opt/logs --name $container_name apollo/apollo-portal:${apollo_version}

    fi
}
