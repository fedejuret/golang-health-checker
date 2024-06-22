# Health Checker

<hr>
Health checker is a tool developed in Golang that allows you to check the status of the web services you want
Supports notifications to different channels to alert you about unwanted https codes.

## Features

- Multiple services
- Loggers per service. And multiple loggers per service
- Timeout configuration
- Http status codes verifications

## How to

Each service is configured in **.json** format within `services/`. You will find an example `service.example.json` which
you can copy to start configuring yours.

Once you have your services configured, simply run the binary depending on your operating system. This starts a
background process and verifies each service based on the "every" that has been configured.

## Run with docker

To run this with Docker, just type:

```bash
docker compose up -d --build
```

## Contact

Get in touch with me!

Email: [fedejuret@gmail.com](mailto:fedejuret@gmail.com)<br>
Website: [https://federicojuretich.com](https://federicojuretich.com)<br>
LinkedIn: [https://www.linkedin.com/in/federicojuretich/](https://www.linkedin.com/in/federicojuretich/)