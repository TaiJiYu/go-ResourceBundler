# go-ResourceBundler
## 1. Usage
### Package resources
1. Run sh.\build.sh
2. Generate the tool directory and store your own resource files under tool/resource. The directory can contain folders
3. Run the packer.exe file in the tool directory to automatically package the resource files
   - You can also run packer.exe from the command line with the following optional parameters:
     -   -n string
     Product name, example: packer.exe -n my_package
     -  -o string
     Product directory, example: packer.exe-o./outcome
     -  -r string
     resource directory, for example, packer.exe -r./resource
     -  -s string
The password must be no less than 16 bytes, for example, packer.exe -s codesec
1. Generate the fqm file in the outcome directory
2. Call.fqm using the following methods:
    ```go
    option := component.UnpackOption{
        FqmFilePath: "tool/outcome/my_resource.fqm", // can provide an absolute or relative directory
        SecretKey: []byte("hjaslkdh"), // password
    }
    err := api.UnpackerInit(o)
        if err ! = nil {
        panic(err)
    }
    api.Show() // View basic information about the fqm file
    fmt.Println(api.Key("filepath/xx.txt"))
    api.Close() // Called when finished
    ```
## 2. File Structure

| No. | Example | Meaning | Number of bytes |
|:-:|:-:|:-:|:-:|
|1|66 71 6D| header |3|
|2|00 00 00 00| Performs the CRC-32 check on all encrypted bytes |4|
|3|01 00| Minimum sdk version 1.0|2|
|4|01| Whether to encrypt,01 encrypted,00 not encrypted |1|
|5|07 E8| year the file was created |2|
|6|0B 08| Months [1] Dates [1]|2|
|7|00 00 00 00| Total number of bytes of index information |4|
|8|00 00 00 00 00 00| Total number of bytes of the index name |6|
|9|00 00 00 00 00 00| Total bytes of data |6|
|10*i|00 00 00 00|key_i start byte position |4|
|10*i+1|00 00| The number of bytes of key_i |2|
|10*i+2|00 00 00 00 00 00 00|key_i corresponds to the start byte of data |6|
|10*i+3|00 00 00 00 00 00 00| Number of bytes of data corresponding to key_i |6|
||** ** ** ** **|The arrangement of key ||
||** ** ** ** **| Arranged by ||