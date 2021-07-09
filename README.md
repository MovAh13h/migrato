# Migrato
Dead simple MySQL Migrations for GO

## Usage

 - Initialize
	```bash
	$ go run https://github.com/movah13h/migrato init 
	```

 - Create new migration
	```bash
	$ go run https://github.com/movah13h/migrato new -name=create_table_auth
    ```

 - Run a specific migration
	```bash
	$ go run https://github.com/movah13h/migrato run -dsn=root:root@tcp(127.0.0.1:3306)/dbName -file=./migrations/up/16101112_create_table_auth.up.sql
    ```

 - Run all `up` migrations
	```bash
	$ go run https://github.com/movah13h/migrato up -dsn=root:root@tcp(127.0.0.1:3306)/dbName
    ```

 - Run all `down` migrations
 	```bash
	$ go run https://github.com/movah13h/migrato down -dsn=root:root@tcp(127.0.0.1:3306)/dbName
    ```