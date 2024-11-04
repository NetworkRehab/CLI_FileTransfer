# CLI File Transfer Utility

⚠️ **WARNING: This project is in pre-alpha stage and not ready for production use** ⚠️

## Description

A command-line utility for transferring files using various protocols such as Azure Blob Storage, Amazon S3, SFTP, CIFS, and local file transfers. The project aims to provide a simple, unified interface for file transfers across different protocols.

## Current Status

This project is under active development and many features are incomplete or not fully tested.

### Working Features
- Basic TUI (Terminal User Interface) with protocol selection
- Configuration file support
- Local file transfer implementation
- Basic database logging of transfers
- Protocol framework for multiple transfer types

### In Progress
- SFTP implementation (partial)
- Azure Blob Storage implementation (partial)
- S3 implementation (partial)
- CIFS/SMB implementation (not started)
- Error handling improvements
- Progress reporting
- Transfer validation
- Credential management
- Unit tests

### Known Issues
- Authentication not fully implemented for remote protocols
- No progress indication during transfers
- Error handling needs improvement
- Configuration validation missing
- No retry mechanism for failed transfers
- Missing proper logging system
- Security considerations need review

## Supported Protocols (Planned)

- Azure Blob Storage (`azureblob`) - Partial
- CIFS/SMB (`cifs`) - Not implemented
- SFTP (`sftp`) - Partial
- Amazon S3 (`s3`) - Partial
- Local file system (`local`) - Working

## Installation

## Usage
Usage
Interactive Mode (Recommended)
This will start the Terminal User Interface (TUI) where you can:

Select a transfer protocol using arrow keys
Enter source path
Enter destination path
Command Line Mode
./CLI_FileTransfer -protocol <protocol> -source <source_path> -destination <dest_path>

# Examples:
# Local file copy
./CLI_FileTransfer -protocol local -source ./myfile.txt -destination ./backup/myfile.txt

# SFTP transfer
./CLI_FileTransfer -protocol sftp -source ./localfile.txt -destination /remote/path/file.txt

# S3 upload
./CLI_FileTransfer -protocol s3 -source ./myfile.txt -destination bucket/myfile.txt

# Azure Blob transfer
./CLI_FileTransfer -protocol azureblob -source ./myfile.txt -destination container/myfile.txt
