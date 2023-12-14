# Tokopedia Crawler

Tokopedia Crawler is a tool for scraping data from tokopedia.com and storing it in a PostgreSQL database or save it to csv file.

## Requirements

Before running the program, please ensure that you have the following dependencies installed:

- Python 3
- Go 1.20
- Chrome browser (for website scraping)
- PostgreSQL

## Installation

### 1. Clone the repository:

```shell

git clone https://github.com/your/repository.git

```
### 2. Install the required dependencies:

For Go dependencies, run the following command in the project root directory:

```shell
go get all
```

For Python dependencies, navigate to the ./other directory and run:

```shell
cd other && pip install -r requirements.txt
```

### 3. Configure the PostgreSQL database:

Create a new PostgreSQL database.
Update the database connection details in the program's configuration file.
Usage

### 4. Run the program:

```
go run main.go -h
```

Follow the prompts and provide the necessary inputs when prompted.

## Option to run

### 1. Scrapping Tokopedia categories https://www.tokopedia.com/p/handphone-tablet/handphone

If you need to test the scrapper don't forget migrate the database first.
```shell
go run main.go scrapper -c=p/handphone-tablet/handphone -t=100 --thread=5 --db --file
```

For other options :
```shell
Use to scrapping data from website

Usage:
  cmd scrapper [flags]

Flags:
  -c, --category string   Category name for scraping
      --db                Save data to DB
      --file              Save data to File
  -h, --help              help for scrapper
      --thread uint       MultiThread scraping (default 5)
  -t, --total uint        Total of data (default 100)
```

### 2. Database Migration
For database migrations, if you dont want to save data, don't use flag *--db* while scraping.

```shell
go run main.go goose up
```
For other options :

```shell
Use to migrate table to database

Usage:
  cmd goose [flags]

Flags:
      --allow-missing   applies missing (out-of-order) migrations
  -d, --dir string      migrations directory (default "src/migration")
  -h, --help            print help
      --no-versioning   apply migration commands with no versioning, in file order, from directory pointed to
      --schema string   for type (migration) or (seeder) (default "migration")
  -s, --sequential      use sequential numbering for new migrations
  -t, --table string    migrations table name (default "goose_db_version")
  -v, --verbose         enable verbose mode
      --version         print version
```