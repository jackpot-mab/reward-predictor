# reward-predictor

- Service that returns predicted and estimated values for each arm in the jackpot's decision-making process.
- The models ONNX files should be saved in a bucket in S3.
- A Go wrapper around onnxruntime is used to predict.
- Some algorithms may need other estimations to be performed apart from the model prediction. Check de docs
for more info.

## Set up

### Set up S3 locally
Set up minio to mock s3:
```
docker run -d --name s3 -p 9000:9000 -p 9001:9001 quay.io/minio/minio server /data --console-address ":9001"
```
- The default user and password is minioadmin.
- Create create credentials and store them in a .env file 
- Create a bucket named "model" and store the models there.

### Set up onnxruntime
The models are supposed to be saved in ONNX format. To load and run inferences on these models, the 
onnxruntime library is used. To facilitate this, we utilize a wrapper for the C onnxruntime API, available 
at https://github.com/yalue/onnxruntime_go. This wrapper requires loading the library files locally.

To acquire the system library files, you can visit https://github.com/microsoft/onnxruntime/releases and 
download the appropriate files. Once downloaded, update the relevant variable in your .env file to point to 
the local location of these library files.


