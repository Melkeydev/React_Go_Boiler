## React/Go Boiler

### Setup

in `main.go` configure new SMTP credentials if you want that functionality.
if not remove.

Modify the parameters in the `docker-compose.yml` file for what you want to name the database and connections

Run `docker-compose up`
This will create an instance of a postgres DN in docker

If successful, the following message should be up:

```
{"level":"INFO","time":"2022-01-09T03:11:42Z","message":"Loading server..."}
{"level":"INFO","time":"2022-01-09T03:11:42Z","message":"Server running on port"}
```
