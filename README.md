# gator
A Blog Aggregator, made for the Boot.dev Backend Learning Path

## Requirements

- PostgreSQL
	- macOS with **brew**:   
	`brew install postgresql@17`  
	- Linux/WSL:  
	`sudo apt update`  
	`sudo apt install postgresql postgresql-contrib`  
- Go  
	- You can download it from the [official download page](https://go.dev/doc/install)

## Installation

To install gator, execute this in your terminal:  
`go install github.com/luism2302/gator/cmd/gator@latest`

Before using gator, you'll need a .gatorconfig.json file in your home directory. Create the file with the following json:  
```
{
	"db_url" : "postgres://username@localhost:5432/database?sslmode=disble"
}
```
Replace url with your own database connection string

## Usage

To create a new user:  
`gator register <username>`  

To login as a previously created user:  
`gator login <username>`  

To see all the registered users:  
`gator users`  

To add a feed as the current user:  
`gator addfeed <feed_name> <feed_url>`  

To see all the added feeds:  
`gator feeds`  

To follow a feed, as the current user:  
`gator follow <feed_url>`  

To unfolow:  
`gator unfollow <feed_url>`  

To see all the followed feeds for the current user:  
`gator following`  

To start the aggregator:  
`gator agg 15s`  

To browse posts, if no limit is provided it defaults to 2:  
`gator browse <limit>`   

