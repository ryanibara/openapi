package openapi_test

// func TestLink(t *testing.T) {
// 	assert := require.New(t)
// 	j := []string{
// 		`{
// 			"address": {
// 			  "operationId": "getUserAddress",
// 			  "parameters": {
// 				"userId": "$request.path.id"
// 			  }
// 			}
// 		  }`,
// 		`{
// 			"UserRepositories": {
// 			  "operationRef": "https://na2.gigantic-server.com/#/paths/~12.0~1repositories~1{username}/get",
// 			  "parameters": {
// 				"username": "$response.body#/username"
// 			  }
// 			}
// 		  }`,
// 	}

// 	for _, d := range j {
// 		data := []byte(d)
// 		var ll openapi.LinkMap
// 		err := json.Unmarshal(data, &ll)
// 		assert.NoError(err)
// 		b, err := json.Marshal(ll)
// 		assert.NoError(err)
// 		assert.True(jsonpatch.Equal(data, b))

// 		// testing yaml

// 		y, err := yaml.JSONToYAML(data)
// 		assert.NoError(err)
// 		var yo openapi.LinkMap
// 		err = yaml.Unmarshal(y, &yo)
// 		assert.NoError(err)
// 		yb, err := json.MarshalIndent(yo, "", "  ")
// 		assert.NoError(err)
// 		if !jsonpatch.Equal(data, yb) {
// 			fmt.Println(string(data), "\n------------------------\n", string(yb))
// 		}
// 		assert.True(jsonpatch.Equal(data, yb), cmpjson.Diff(data, yb))

// 	}
// }
