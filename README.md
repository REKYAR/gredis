# DOCS

## RUNNING SERVER

## RUNNING CLIENT

## PROTOCOL DOCUMENTATION (RECIEVE)
message structure in order  
**protocol version (uint64)**  
rest of message, determined by protocol version

 - ### version (1)
     - assumes single CRUD command + data

    **total size = message size + metadata size (bytes) (uint64)**   
    **metadata size (bytes) (uint64)**  
    **metadata**   
    **data**  

## PROTOCOL DOCUMETATION (RETURN)
message structure in order  
**protocol version (uint64) - the same as recieved**  
rest of message, determined by protocol version
  - ### version (1)
    - produced after completion/failure of the whole read/write

    **total size = message size + metadata size (bytes) (uint64)**   
    **metadata size (bytes) (uint64)**  
     - contains operation status under key "status"
     - values are OK for success and FAIL for failure

    **message**
     - if operation failed you will find a string message here, if operation was successful it will be empty


#### Notes to self
- net.Conn.Read does not convert endiannes