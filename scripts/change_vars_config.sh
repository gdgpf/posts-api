cd $1
ls
cp config.example.json config.json
sed -i -e 's/$BCRYPT_HASH/'dev'/g' config.json
sed -i -e 's/$DB_USER/'postgres'/g' config.json
sed -i -e 's/$DB_PASS/'postgres'/g' config.json
sed -i -e 's/$DB_HOST/'postgres'/g' config.json
sed -i -e 's/$DB_PORT/'5432'/g' config.json
sed -i -e 's/$DB_NAME/'dbdev'/g' config.json
sed -i -e 's/$SERVER_PORT/'3000'/g' config.json
sed -i -e 's/$SERVER_VERSION/'2022.6.5.0'/g' config.json