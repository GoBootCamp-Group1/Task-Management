#!/bin/bash
set -e

# Check if database exists
psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'task-management'" | grep -q 1 || psql -U postgres -c "CREATE DATABASE \"task-management\""

# Create tables
psql -U postgres -d "task-management" -f /app/Task-manager.sql