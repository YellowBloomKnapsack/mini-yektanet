# Mini Yektanet

## What is It?
Mini [Yektanet](https://yektanet.com) is an imitation of an online advertising platform that suggests and shows ads in publishers websites and keeps tracks of clicks/impressions.
The project was done as a part of Yektanet's Yellow Bloom workshop, in two weeks.
The app is made of several components:
1. **Ad Server**: Responsible for giving ads to publishers, whenever wanted.
2. **Event Server**: Handles click and impression events.
3. **Panel**: Has dashboards for customers and is the interface to the database.
4. **Publishers Website**: A few number of  hard-coded web pages.

## How to run

You can use docker to build and run the project. However if you prefer to run it locally, here are the instructions:

### Requirements
You'll need to have these tools installed and running on your system:
1. Kafka
2. PostgreSQL
3. Redis (with Bloom Filter module)

The configuration on how to connect to external components are in the environment variables.

And of course, you'll need Golang and the Golang packages used in the project (can be installed using `go`) 


### Environment Variables
Copy the contents of `.env.example` to `.env` in the same directory.

**Note**: Some values are secrets and not included in the `.env.example`. Like database password and private keys.

### Running Golang Apps
Navigate to each service's directory and run each one of them:
```
cd adserver
go run .
```
The list of services you need to run:
1. Ad Server (`./adserver/`)
2. Event Server (`./eventserver/`)
3. Panel (`./panel/`)
4. Publishers Websites (`./publisherwebsite/`)

Now go to localhost:8084/varzesh3
Sorry if there is no fancy UI (or no UI at all), we did not have time for that.

## Contributors

Developers:
- [Amir Parsa Khadem](https://github.com/aparsak)
- [Nima Rezaei](https://github.com/Haximos)
- [Mohammad Sadegh Poulaei](https://github.com/MSPoulaei)
- [Kourosh Alinaghi](https://github.com/KouroshAlinaghi)

Mentor:
- [Javad Karimi](https://github.com/JKarimi12)

By the way, the name of our team was knapsack. Can you tell why?
