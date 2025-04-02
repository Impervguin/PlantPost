mc alias set minio http://$MINIO_HOST:$MINIO_PORT $MINIO_ROOT_USER $MINIO_ROOT_PASSWORD
mc mb --ignore-existing minio/"$PLANT_BUCKET"
mc anonymous set public minio/"$PLANT_BUCKET"
mc mb --ignore-existing minio/"$POST_BUCKET"
mc anonymous set public minio/"$POST_BUCKET"