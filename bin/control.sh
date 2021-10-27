#!/usr/bin/env bash

. /etc/profile

YEAR_WEEK=$(date +%G%V)

APP_NAME="facm"
APP_BIN="${APP_NAME}"

idc=$(cat /data/wwwroot/idc.dat)

# 先检查状态，进程正常
function start() {
    staus >/dev/null
    staus=$?

    if  [[  ${status} = 1 ]];  then
    echo  "The ${APP_NAME} is already running"
    return 0
    fi
    if [[  ${status} = 2 ]];  then
        echo "The ${APP_NAME} is start error"
        return 0
    fi

    NOW=`date "+%Y%m%d"`
    LOG_FILE=$LOG_PATH"/FacmServer_"${NOW}".out"


}