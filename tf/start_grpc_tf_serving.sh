docker run -p 8500:8500 \
--mount type=bind,source=/Users/sauganth/go/src/github.com/sauganth/goprojects/tf/model,target=/models/dense \
-e MODEL_NAME=dense -t tensorflow/serving &
