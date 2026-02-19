# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This repository is for designing and implementing a Web server in Go which has the following features

- Index or homepage: A simple Single page application using bootstrap5 templates
- REST API: Server at /api/v1/ with endpoints for CRUD operations on a resource (e.g., users, products), in-memory data storage, and JSON request/response handling
- WebSocket: Real-time communication endpoint at /ws/ for broadcasting messages to connected clients

## IMPORTANT:

- This project uses LiteLLM. Always make one tool call at a time. Never use more than one tool call. This is to ensure that the system remains responsive and efficient while processing requests.
