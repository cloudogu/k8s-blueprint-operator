package blueprint

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"

	k8sv1 "github.com/cloudogu/k8s-blueprint-operator/pkg/api/v1"
)

var testCtx = context.Background()

func Test_blueprintClient_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, "GET", request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/blueprints/testblueprint", request.URL.Path)
			assert.Equal(t, http.NoBody, request.Body)

			writer.Header().Add("content-type", "application/json")
			blueprint := &k8sv1.Blueprint{ObjectMeta: v1.ObjectMeta{Name: "testblueprint", Namespace: "test"}}
			blueprintBytes, err := json.Marshal(blueprint)
			require.NoError(t, err)
			_, err = writer.Write(blueprintBytes)
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := newForConfig(&config)
		require.NoError(t, err)
		dClient := client.Blueprints("test")

		// when
		_, err = dClient.Get(testCtx, "testblueprint", v1.GetOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_blueprintClient_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodGet, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/blueprints", request.URL.Path)
			assert.Equal(t, http.NoBody, request.Body)

			writer.Header().Add("content-type", "application/json")
			blueprintList := k8sv1.BlueprintList{}
			blueprint := &k8sv1.Blueprint{ObjectMeta: v1.ObjectMeta{Name: "testblueprint", Namespace: "test"}}
			blueprintList.Items = append(blueprintList.Items, *blueprint)
			blueprintBytes, err := json.Marshal(blueprintList)
			require.NoError(t, err)
			_, err = writer.Write(blueprintBytes)
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := newForConfig(&config)
		require.NoError(t, err)
		dClient := client.Blueprints("test")

		// when
		_, err = dClient.List(testCtx, v1.ListOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_blueprintClient_Watch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, "GET", request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/blueprints", request.URL.Path)
			assert.Equal(t, http.NoBody, request.Body)
			assert.Equal(t, "labelSelector=test&watch=true", request.URL.RawQuery)

			writer.Header().Add("content-type", "application/json")
			_, err := writer.Write([]byte("egal"))
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := newForConfig(&config)
		require.NoError(t, err)
		dClient := client.Blueprints("test")

		// when
		_, err = dClient.Watch(testCtx, v1.ListOptions{LabelSelector: "test"})

		// then
		require.NoError(t, err)
	})
}

func Test_blueprintClient_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		blueprint := &k8sv1.Blueprint{ObjectMeta: v1.ObjectMeta{Name: "tocreate", Namespace: "test"}}

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPost, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/blueprints", request.URL.Path)

			bytes, err := io.ReadAll(request.Body)
			require.NoError(t, err)

			createdBlueprint := &k8sv1.Blueprint{}
			require.NoError(t, json.Unmarshal(bytes, createdBlueprint))
			assert.Equal(t, "tocreate", createdBlueprint.Name)

			writer.Header().Add("content-type", "application/json")
			_, err = writer.Write(bytes)
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := newForConfig(&config)
		require.NoError(t, err)
		dClient := client.Blueprints("test")

		// when
		_, err = dClient.Create(testCtx, blueprint, v1.CreateOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_blueprintClient_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		blueprint := &k8sv1.Blueprint{ObjectMeta: v1.ObjectMeta{Name: "tocreate", Namespace: "test"}}

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPut, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/blueprints/tocreate", request.URL.Path)

			bytes, err := io.ReadAll(request.Body)
			require.NoError(t, err)

			createdBlueprint := &k8sv1.Blueprint{}
			require.NoError(t, json.Unmarshal(bytes, createdBlueprint))
			assert.Equal(t, "tocreate", createdBlueprint.Name)

			writer.Header().Add("content-type", "application/json")
			_, err = writer.Write(bytes)
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := newForConfig(&config)
		require.NoError(t, err)
		dClient := client.Blueprints("test")

		// when
		_, err = dClient.Update(testCtx, blueprint, v1.UpdateOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_blueprintClient_UpdateStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		blueprint := &k8sv1.Blueprint{ObjectMeta: v1.ObjectMeta{Name: "tocreate", Namespace: "test"}}

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPut, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/blueprints/tocreate/status", request.URL.Path)

			bytes, err := io.ReadAll(request.Body)
			require.NoError(t, err)

			createdBlueprint := &k8sv1.Blueprint{}
			require.NoError(t, json.Unmarshal(bytes, createdBlueprint))
			assert.Equal(t, "tocreate", createdBlueprint.Name)

			writer.Header().Add("content-type", "application/json")
			_, err = writer.Write(bytes)
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := newForConfig(&config)
		require.NoError(t, err)
		dClient := client.Blueprints("test")

		// when
		_, err = dClient.UpdateStatus(testCtx, blueprint, v1.UpdateOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_blueprintClient_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodDelete, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/blueprints/testblueprint", request.URL.Path)

			writer.Header().Add("content-type", "application/json")
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := newForConfig(&config)
		require.NoError(t, err)
		dClient := client.Blueprints("test")

		// when
		err = dClient.Delete(testCtx, "testblueprint", v1.DeleteOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_blueprintClient_DeleteCollection(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodDelete, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/blueprints", request.URL.Path)
			assert.Equal(t, "labelSelector=test", request.URL.RawQuery)
			writer.Header().Add("content-type", "application/json")
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := newForConfig(&config)
		require.NoError(t, err)
		dClient := client.Blueprints("test")

		// when
		err = dClient.DeleteCollection(testCtx, v1.DeleteOptions{}, v1.ListOptions{LabelSelector: "test"})

		// then
		require.NoError(t, err)
	})
}

func Test_blueprintClient_Patch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPatch, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/blueprints/testblueprint", request.URL.Path)
			bytes, err := io.ReadAll(request.Body)
			require.NoError(t, err)
			assert.Equal(t, []byte("test"), bytes)
			result, err := json.Marshal(k8sv1.Blueprint{})
			require.NoError(t, err)

			writer.Header().Add("content-type", "application/json")
			_, err = writer.Write(result)
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := newForConfig(&config)
		require.NoError(t, err)
		dClient := client.Blueprints("test")

		patchData := []byte("test")

		// when
		_, err = dClient.Patch(testCtx, "testblueprint", types.JSONPatchType, patchData, v1.PatchOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_blueprintClient_UpdateStatusXXX(t *testing.T) {
	for _, testCase := range []struct {
		functionName   string
		expectedStatus k8sv1.StatusPhase
	}{
		{
			functionName:   "UpdateStatusInProgress",
			expectedStatus: k8sv1.StatusPhase("inProgress"),
		},
		{
			functionName:   "UpdateStatusCompleted",
			expectedStatus: k8sv1.StatusPhase("completed"),
		},
		{
			functionName:   "UpdateStatusInvalid",
			expectedStatus: k8sv1.StatusPhase("invalid"),
		},
		{
			functionName:   "UpdateStatusFailed",
			expectedStatus: k8sv1.StatusPhase("failed"),
		},
		{
			functionName:   "UpdateStatusRetrying",
			expectedStatus: k8sv1.StatusPhase("retrying"),
		},
	} {
		t.Run(fmt.Sprintf("%s success", testCase.functionName), func(t *testing.T) {
			// given
			blueprint := &k8sv1.Blueprint{ObjectMeta: v1.ObjectMeta{Name: "testblueprint", Namespace: "test"}}

			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				switch request.Method {
				case http.MethodGet:
					assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/blueprints/testblueprint", request.URL.Path)
					assert.Equal(t, http.NoBody, request.Body)

					writer.Header().Add("content-type", "application/json")
					blueprint := &k8sv1.Blueprint{ObjectMeta: v1.ObjectMeta{Name: "testblueprint", Namespace: "test"}}
					blueprintBytes, err := json.Marshal(blueprint)
					require.NoError(t, err)
					_, err = writer.Write(blueprintBytes)
					require.NoError(t, err)
					writer.WriteHeader(200)
				case http.MethodPut:
					assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/blueprints/testblueprint/status", request.URL.Path)
					bytes, err := io.ReadAll(request.Body)
					require.NoError(t, err)

					createdBlueprint := &k8sv1.Blueprint{}
					require.NoError(t, json.Unmarshal(bytes, createdBlueprint))
					assert.Equal(t, "testblueprint", createdBlueprint.Name)
					assert.Equal(t, testCase.expectedStatus, createdBlueprint.Status.Phase)

					writer.Header().Add("content-type", "application/json")
					_, err = writer.Write(bytes)
					require.NoError(t, err)
					writer.WriteHeader(200)
				default:
					assert.Fail(t, "method should be get or put")
				}
			}))

			config := rest.Config{
				Host: server.URL,
			}
			client, err := newForConfig(&config)
			require.NoError(t, err)
			dClient := client.Blueprints("test")

			// when
			returnValues := reflect.ValueOf(dClient).MethodByName(testCase.functionName).Call([]reflect.Value{reflect.ValueOf(testCtx), reflect.ValueOf(blueprint)})
			err, _ = returnValues[1].Interface().(error)

			// then
			require.NoError(t, err)
		})
	}
}
