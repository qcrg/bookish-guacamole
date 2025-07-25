#!/bin/env bash
set -e

TEMP_FILE=`mktemp`

if ! ./health.sh; then
  echo Health error 1>&2
fi

export USER_ID="01983784-e5ff-786a-bdd0-ab93f98362d2"
echo USER_ID=$USER_ID

export VERSION=`./version.sh`
echo VERSION=$VERSION

./auth.sh > ${TEMP_FILE}
source ${TEMP_FILE}
echo ACCESS_TOKEN=$ACCESS_TOKEN
echo REFRESH_TOKEN=$REFRESH_TOKEN

ME_USER_ID=`./me.sh`
if [ 0 -ne $? ]; then
  echo me user_id != user_id
fi

echo $TEMP_FILE

./refresh_session.sh > ${TEMP_FILE}
source ${TEMP_FILE}
./refresh_session.sh > ${TEMP_FILE}
source ${TEMP_FILE}
./refresh_session.sh > ${TEMP_FILE}
source ${TEMP_FILE}
./refresh_session.sh > ${TEMP_FILE}
source ${TEMP_FILE}
./refresh_session.sh > ${TEMP_FILE}
source ${TEMP_FILE}

./delete_session.sh

rm ${TEMP_FILE}
