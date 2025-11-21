# GoCronus

## Overview
**GoCronus** is a CRON job notification system designed to schedule and send notifications. It supports scheduling notifications across **CST** and **IST** timezones, making it suitable for applications with global user bases. The project is developed in **Go (Golang)**, using **MySQL** as the database for storing notification details and **Google Firebase Cloud Messaging (FCM)** for sending notifications.

The project leverages **sqlboiler** as the Object Relational Mapping (ORM) tool to interact with the database. It includes multiple services working together to handle scheduling, sending, and managing notifications efficiently. Logging is implemented using **zap** (for structured logging) and **lumberjack** (for log rotation), ensuring high-performance logging with effective file management.

---

## Features
- **Timezone Support:** Schedule notifications in both CST and IST timezones.  
- **Efficient CRON Jobs:** Robust job scheduling and execution.  
- **MySQL Database:** Stores notification details, user information, and scheduling data.  
- **Google FCM Integration:** Delivers push notifications to client devices.  
- **ORM Integration:** Uses sqlboiler for clean and efficient database interactions.  
- **Modular Codebase:** Services are decoupled for better maintainability and scalability.  
- **High-Performance Logging:** Utilizes zap for structured logging and lumberjack for log rotation.  
- **Adherence to Best Practices:** Implements clean code principles, proper error handling, and structured project organization.  

---

## Tech Stack
- **Language:** Go (Golang)  
- **Database:** MySQL  
- **Messaging:** Google Firebase Cloud Messaging (FCM)  
- **ORM:** sqlboiler  
- **Job Scheduler:** CRON jobs managed within Go services  
- **Logging:** zap (structured logging) and lumberjack (log rotation)  
- **Deployment:** Docker and Kubernetes support (if applicable)  

---

## Project Structure
```plaintext
gocronus/
├── cmd/                   # Main entry points for different services
│   ├── notifier/          # Service to send notifications
│   └── scheduler/         # CRON job service for scheduling notifications
├── models/                # Auto-generated ORM models using sqlboiler
├── configs/               # Configuration files (e.g., database, FCM, CRON settings)
├── services/              # Business logic for scheduling and sending notifications
├── utils/                 # Helper functions and common utilities
├── pkg/logger/            # Logging setup using zap and lumberjack
├── logs/                  # Log files for monitoring and debugging
└── docker-compose.yml     # Docker configuration for local development
```

---

## How It Works
- **Notification Scheduling:**  
  - Users or services request notification scheduling through an API.  
  - The **scheduler service** sets up CRON jobs based on the requested time and timezone.  

- **Notification Sending:**  
  - When the scheduled time arrives, the **notifier service** triggers and retrieves notification details from the MySQL database.  
  - Notifications are sent using **Google FCM** to target devices.  

- **Database Interaction:**  
  - Uses **sqlboiler** to map Go structs to MySQL tables for seamless database operations.  

- **Timezone Handling:**  
  - Scheduling accounts for CST and IST timezones to ensure timely notifications for users across regions.  

- **Logging:**  
  - Logging is handled by **zap** for high-performance structured logging.  
  - **lumberjack** is integrated for log rotation, ensuring logs don’t consume excessive disk space.  
  - Logs include details like timestamps, log levels, file paths, and messages.  

---

## Best Practices Followed
- **Structured Logging:** Provides consistent, machine-readable logs using zap.  
- **Log Rotation:** Prevents excessive disk usage through lumberjack’s log management.  
- **Separation of Concerns:** Distinct services for scheduling and notification dispatching.  
- **Environment Variables:** Secure handling of sensitive data through `.env` files.  
- **Error Handling:** Clear error messages with detailed context for faster debugging.  
- **Resource Management:** Ensures proper closing of database connections and CRON jobs.  
- **Security:** Excludes sensitive data like passwords from logs.  

---

## Future Improvements
- Web dashboard for managing and monitoring notifications  
- Enhanced logging with request tracing and correlation IDs  
- Retry mechanism for failed notifications  
- Multi-timezone and multi-language notification support  