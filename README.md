Smart Calculator

A full-stack web application with a Go backend and an HTML/JS frontend. This project was built to learn how to connect a local Go API to a live web interface hosted on GitHub Pages and Render.
Features

    Full-Stack Integration: Uses a Go REST API for all mathematical logic.

    Modern UI: Clean interface with CSS Grid and glassmorphism effects.

    Live Status: Real-time server health check to monitor the backend connection.

    Deployment: Frontend hosted on GitHub Pages; Backend hosted on Render.

Project files

    .gitignore: Lists files and directories to be ignored by Git.

    go.mod: The Go module definition file that manages project dependencies.

    index.html: The frontend user interface and JavaScript logic for API communication.

    README.md: Documentation providing an overview and setup instructions for the project.

    main.go: The backend server containing API routes and mathematical calculation logic.

Setup and Development
Backend

The backend is written in Go. It handles the calculations and manages CORS to allow requests from the GitHub Pages frontend.

    Initialize and Run:
    Bash

    go mod init calculator
    go run main.go

    Key Endpoint: POST /calculate – Receives a JSON expression and returns the result.

Frontend

The frontend is a single-page application (SPA).

    API_URL: Points to the live Render URL https://kalkutor.onrender.com/calculate.

    Deployment: The index.html file is placed in the root directory for automatic deployment via GitHub Pages.

API Endpoints

    GET /health: Returns service status.

    POST /calculate: Accepts {"expression": "string"} and returns the computed result.

Deployment Notes

    Frontend: Hosted at https://eheguy.github.io/kalkutor/

    Backend: Hosted on Render. Note that the free tier may take 30–60 seconds to spin up after inactivity.

    CORS: The Go backend includes headers to allow the GitHub domain to access the API.

Lessons Learned

    Handling Cross-Origin Resource Sharing (CORS) in Go.

    Managing environment variables like PORT for cloud deployment.

    Troubleshooting Git push conflicts and repository structure for GitHub Pages.





  ~i wrote the readme through the help of chat gpt its late now~ im going to sleep any suggesitons are welcone~
