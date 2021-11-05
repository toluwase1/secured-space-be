# Rentals-API

Rentals-API is a `RESTful-API` that allows `developers` to ` connect to data about rental properties`.


## Using Rentals-API

To use Rentals-API, clone this repository and follow one of two methods shown below:

### Using Docker Compose

1. [Install Docker Compose](https://docs.docker.com/compose/install/)
2. Run all containers with `make up`.


### Using Your Local Machine

1. [Install Go](https://golang.org/doc/install)
2. [Install PostgreSQL](https://www.postgresql.org/download/)
3. Create a database named `rental`.
4. Install all dependencies using:


```
go mod tidy
```
5. Run the application using:

```
make run
```

## Contributing to Rentals-API

To contribute to Rentals-API, follow these steps:

1. Clone this repository.
2. Create a branch: `git checkout -b <branch_name>`. When you create a branch, it should follow the name of the task.
3. Make your changes and commit them: `git commit -m '<commit_message>'`
4. Push to the original branch: `git push origin <branch_name>/<location>`
5. Create the pull request.

Alternatively see the GitHub documentation on [creating a pull request](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request).

## Using the Rentals-API Email Template
Mailchimp api  is used for sending emails to customers using the platform.

1. To gain access to this template, log on to ([mailchimp](https://login.mailchimp.com))  you need to login with the rental-api verified username and password.
2. Navigate to the create campaigns, then click on email templates.
3. Click on the rental-api template.
4. Edit the template, add the links desired and save the template.

## Setting Environmental Variables
An environment variable is a text file containing ``KEY=value`` pairs of your secret keys and other private information. For security purposes, it is ignored using ``.gitignore`` and not committed with the rest of your codebase.

To create, ensure you are in the root directory of the project then on your terminal type:
```
touch .env
```
All the variables used within the project can now be added within the ``.env`` file in the following format:
```
DB_HOST=127.0.0.1
DB_PORT=8080
DB_USER=<your db username>
DB_PASS=<your db password>
```


## Tests
Testing is done using the GoMock framework. The ``gomock`` package and the ``mockgen``code generation tool are used for this purpose.
If you installed the dependencies using the command given above, then the packages would have been installed. Otherwise, installation can be done using the following commands:
```
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen
```

After installing the packages, run:
```
make mock-db
```

The command above helps to generate mock interfaces from a source file.

To run tests, run:
```
make test
```

## Contributors

Thanks to the following people who have contributed to this project:

* Odohi David ([spankie](https://github.com/spankie)) 📖
* Toluwase Thomas ([toluwase1](https://github.com/toluwase1)) 📖
* Olusola Alao ([olusolaa](https://github.com/olusolaa)) 📖
* Tambarie Gbaragbo ([Tambarie](https://github.com/Tambarie)) 🐛
* Clinton Adebayo ([Ad3bay0c](https://github.com/Ad3bay0c)) 🐛
* Omoyemi Arigbanla ([yemmyharry](https://github.com/yemmyharry)) 🐛
* Victor Anyimukwu ([udodinho](https://github.com/udodinho)) 🐛
* Franklyn Omonade ([nade-harlow](https://github.com/nade-harlow)) 🐛
* Chisom Amadi ([Tchisom17](https://github.com/Tchisom17)) 🐛
* Nonso Okike ([okikechinonso](https://github.com/okikechinonso)) 🐛
* Shuaib Olurode ([OShuaib](https://github.com/OShuaib)) 🐛
* Chukwuebuka Iroegbu ([iBoBoTi](https://github.com/iBoBoTi)) 🐛
* Nnah Nnamdi ([techagentng](https://github.com/techagentng)) 🐛



