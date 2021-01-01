# File Data Store using Go

File Data store is a CLI application that is capable of storing, retrieving and deleting JSON objects in a CSV file. It has the capability of autodeletion when the expiry date of a certain passes its time.
## Requirements and Specifications

- go version go1.13.8 linux/amd64
- Linux based OS

Other Constraits are:
- The size of the file must not exceed 1GB.
- Must be thread safe
- Key must be below 32 chars and JSON object size must not exceed 16KB.

## Usage
1. Create:
This operation is used to store a JSON object. The input needed are key(unique), JSON, time to live(Optional).
2. Read:
This operation is used to read a JSON object from the file. The input needed is a Key.
3. Delete:
This operation is used to delete a JSON object. The input needed are needed is a Key.
4. Exit:
This operation is used to from the application

## Images
<a href = "https://drive.google.com/file/d/19Bj2qbymOTWm5bqPHemgNqmdeUo5brBd/view?usp=sharing">Run Operation Image</a>

<a href = "https://drive.google.com/file/d/16Hw4mJawbXij0PZs8V8WeFxfqabh0Per/view?usp=sharing">Create Operation Image</a>

<a href = "https://drive.google.com/file/d/1LodafV_CBuhp9mewthIIzA2PT2X31sKQ/view?usp=sharing">Read Operation Image</a>

<a href = "https://drive.google.com/file/d/1ZZTMzLf7evHwA6k06vx6Qi-Pm4XTh3kZ/view?usp=sharing">Before Delete Image</a>

<a href = "https://drive.google.com/file/d/1x8NoOd4WH0wJKrP4eZj-0j-4g1VcRwOQ/view?usp=sharing">After Delete Image</a>

<a href = "https://drive.google.com/file/d/1xSxAHgZhZBmzEQVpRceM50a-_sR6KTTX/view?usp=sharing">Exit Operation Image</a>

<a href = "https://drive.google.com/file/d/1hnJ14CyvlPI2RfcwxWKUfUlPPd_CZne1/view?usp=sharing">CSV File Image</a>


## Issues facing now
I have implemented an expiry date function which deletes the key along with its other data. I am finding a little bit to synchronize it.


## Conclusion
Thanks to Freshworks for this wonderful opportunity and I hope I can learn and work more from there.
