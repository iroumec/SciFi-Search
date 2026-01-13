# Running the Application

The application provides two modes: development mode and production mode.

## Production Mode

This mode consists of a lightweight image, containing only what is strictly necessary to run the application.

To start the application in this mode, run the following command:

```bash
make up
```

The initial startup may take a few minutes. Once completed, the server will run in the background and you will be able to access it at: [http://localhost:8080/](http://localhost:8080/).

To stop the server, run the following command:

```bash
make down
```

This will stop the server and remove the container and associated networks. However, the images will remain and the volumes will persist.

> [!IMPORTANT]
> External volumes are not removed, only those defined within the `docker-compose.yml`.

## Development Mode

This mode has the particularity that, in addition to using an image that includes Golang tools, it integrates Air, which allows changes made to files to be reflected automatically, facilitating development.

To start the application in this mode, run the following command:

```bash
make development
```

In this mode, the server runs in the foreground, so `make down` is not required to stop it. Using `Ctrl + C` is sufficient.

## Cleanup

To remove the image and volumes (useful during development when database changes occur), run the following command:

```sh
make clean
```

**It is also recommended to run this command if you no longer plan to use the application.**
