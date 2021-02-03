# Guest List

A guest list app for the technical test

## Usage
To start the app on port 8080
`docker-compose up`

To run the tests
`make test`

## How does it work?
Its a Golang app that connects to MySQL 5.7

It uses:
- Gorm as a MySQL ORM
- Gorilla Mux as a HTTP Router

There are two data types:
- `Table` which contains a tables `Number` (A chosen designation, to match the real world venue) and its amount of seats
- `Reservation` a booking made for a table covering the whole night. It contains who the guest is, how many 
  accompanying guests they're bringing and which table they have chosen (a foriegn key) to the `Table` type. 
  It doubles as a record of who has arrived by updating the `ArrivalTime`, requiring only one table to store arrived 
  and reserved guests. When this is displayed to the user it is formatted depending on the API using `FormatAsReservation` 
  for the `guest_list` api and `FormatAsGuestArrival` for the `guest` api

The server package starts a web server, which also initialises the `database` that routes to a series of Handlers which 
communicate with the Database package by fetching a singleton by calling `database.Get()`.

### API

#### Tables
You can add and delete tables that are available by providing a table number, and the amount of seats
You can only delete tables if there are no reservations on those tables
```
POST   /table/{number} { "seats": numberOfSeats }
DELETE /table/{number}
GET    /table/{number}
```

#### Reservations
You can reserve tables in advanced of the party, by specifying your name and how many guests you are bringing
```
POST   /guest_list/{name} { "table": tableNumber, "accompanying_guests": numberOfGuests }
DELETE /guest_list/{name}
GET    /guest_list
```

#### Arrivals
When you want to arrive at the party you can check in by providing your name, and how many guests you brought
```
PUT /guest/{name} { "accompanying_guests": numberOfGuests }
GET /guest/{name}
GET /guests
```

#### Other
You can count the amount of empty seats at the party right now, not including people who haven't arrived
```
GET /seats_empty
```

You can also generate HTML invites to the party!
```
GET /inviation/{name}
```


## What could I improve on?
The HTTP response codes are all over the place, with more time I would make sure they are all correct and not just a mix
of bad requests and internal server errors

Also using Gorm cost a lot of time when setting up the project, and the database package is still quite weak setup using a 
singleton and pulling environment variables in a relatively unreliable way, in the future I would have the handlers
be part of a struct and share the DB object which would make it possible for mocking in the future.

The handlers also container nearly all the logic, I would have preferred to have the logic away from the handlers, but
it felt like overkill to separate that right now. Theres lots of room for more helpers to prevent repeat logic,
but I wanted to avoid bloat, because Gorm is already pretty short and easy to handle so adding many more helpers
would have made it harder for a small amount of readability, especially given the scope.

There's also no logging, I would add logurs and implement logging especially on failures.
