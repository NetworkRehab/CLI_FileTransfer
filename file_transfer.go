package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "net/url"
    "os"
    "path/filepath"

    "github.com/Azure/azure-storage-blob-go/azblob"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/pkg/sftp"
    "github.com/spf13/viper"
    "golang.org/x/crypto/ssh"
    "io/ioutil"
)

// transferFile transfers a file using the specified protocol.
// Supported protocols: azureblob, cifs, sftp, s3, local.
func transferFile(protocol, source, destination string) error {
    switch protocol {
    case "azureblob":
        return transferAzureBlob(source, destination)
    case "cifs":
        return transferCIFS(source, destination)
    case "sftp":
        return transferSFTP(source, destination)
    case "s3":
        return transferS3(source, destination)
    case "local":
        return transferLocal(source, destination)
    default:
        return fmt.Errorf("unsupported protocol: %s", protocol)
    }
}

// transferAzureBlob uploads a file to Azure Blob Storage.
func transferAzureBlob(source, destination string) error {
    log.Println("Transferring via Azure Blob Storage")

    accountName := viper.GetString("azure.accountName")
    accountKey := viper.GetString("azure.accountKey")
    containerName := viper.GetString("azure.containerName")

    if accountName == "" || accountKey == "" || containerName == "" {
        return fmt.Errorf("Azure storage credentials not provided")
    }

    credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
    if err != nil {
        return fmt.Errorf("failed to create Azure credential: %v", err)
    }

    pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
    url, _ := url.Parse(
        fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

    containerURL := azblob.NewContainerURL(*url, pipeline)
    blobURL := containerURL.NewBlockBlobURL(filepath.Base(destination))

    file, err := os.Open(source)
    if err != nil {
        return fmt.Errorf("failed to open source file: %v", err)
    }
    defer file.Close()

    ctx := context.Background()
    _, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{})
    if err != nil {
        return fmt.Errorf("failed to upload file to Azure Blob: %v", err)
    }

    log.Println("File transferred via Azure Blob Storage successfully")
    return nil
}

// transferCIFS copies a file to a CIFS (SMB) network share.
func transferCIFS(source, destination string) error {
    log.Println("Transferring via CIFS")

    mountPoint := viper.GetString("cifs.mountPoint")
    if mountPoint == "" {
        return fmt.Errorf("CIFS mount point not configured")
    }

    destPath := filepath.Join(mountPoint, destination)
    return transferLocal(source, destPath)
}

// transferSFTP uploads a file via SFTP.
func transferSFTP(source, destination string) error {
    // Read the allowed host key from a file
    publicKeyBytes, err := ioutil.ReadFile("allowed_hostkey.pub")
    if err != nil {
        return fmt.Errorf("failed to read allowed host key: %v", err)
    }

    // Parse the allowed host key
    publicKey, err := ssh.ParsePublicKey(publicKeyBytes)
    if err != nil {
        return fmt.Errorf("failed to parse allowed host key: %v", err)
    }

    // SSH client configuration
    config := &ssh.ClientConfig{
        User: "username",
        Auth: []ssh.AuthMethod{
            ssh.Password("password"),
        },
        HostKeyCallback: ssh.FixedHostKey(publicKey),
    }

    // Connect to the SSH server
    sshClient, err := ssh.Dial("tcp", "sftp.example.com:22", config)
    if (err != nil) {
        return fmt.Errorf("failed to dial: %v", err)
    }
    defer sshClient.Close()

    // Create new SFTP client
    sftpClient, err := sftp.NewClient(sshClient)
    if (err != nil) {
        return fmt.Errorf("failed to create sftp client: %v", err)
    }
    defer sftpClient.Close()

    // Open source file
    srcFile, err := os.Open(source)
    if (err != nil) {
        return fmt.Errorf("failed to open source file: %v", err)
    }
    defer srcFile.Close()

    // Create destination file on the server
    dstFile, err := sftpClient.Create(destination)
    if (err != nil) {
        return fmt.Errorf("failed to create destination file: %v", err)
    }
    defer dstFile.Close()

    // Copy data from source file to destination file
    _, err = io.Copy(dstFile, srcFile)
    if (err != nil) {
        return fmt.Errorf("failed to copy file: %v", err)
    }

    log.Println("File transferred via SFTP successfully")
    return nil
}

// transferS3 uploads a file to Amazon S3.
func transferS3(source, destination string) error {
    log.Println("Transferring via Amazon S3")

    bucket := viper.GetString("aws.bucket")
    region := viper.GetString("aws.region")

    if bucket == "" || region == "" {
        return fmt.Errorf("AWS S3 bucket or region not configured")
    }

    sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
    if err != nil {
        return fmt.Errorf("failed to create AWS session: %v", err)
    }

    uploader := s3.New(sess)
    file, err := os.Open(source)
    if err != nil {
        return fmt.Errorf("failed to open source file: %v", err)
    }
    defer file.Close()

    _, err = uploader.PutObject(&s3.PutObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(destination),
        Body:   file,
    })
    if err != nil {
        return fmt.Errorf("failed to upload to S3: %v", err)
    }

    log.Println("File transferred via Amazon S3 successfully")
    return nil
}

// transferLocal copies a file locally.
func transferLocal(source, destination string) error {
    // Copying file from source to destination locally
    inputFile, err := os.Open(source)
    if err != nil {
        return fmt.Errorf("failed to open source file: %v", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(destination)
    if err != nil {
        return fmt.Errorf("failed to create destination file: %v", err)
    }
    defer outputFile.Close()

    _, err = io.Copy(outputFile, inputFile)
    if err != nil {
        return fmt.Errorf("failed to copy file: %v", err)
    }

    log.Println("File transferred locally successfully")
    return nil
}
