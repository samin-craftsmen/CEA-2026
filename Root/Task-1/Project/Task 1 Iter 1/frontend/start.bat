#!/bin/bash

echo "Installing frontend dependencies..."
npm install || exit 1

echo "Starting development server..."
npm run dev
