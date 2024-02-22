package config

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestPublicKeyConfigEncryptionAdapter_Encrypt(t *testing.T) {
	t.Run("can encrypt and decrypt", func(t *testing.T) {
		// given
		doguname := "testdogu"
		namespace := "testingnamespace"
		testValue := "value"
		keyProviderStr := "pkcs1v15"
		//encryptedTestValue := common.EncryptedDoguConfigValue("Z5MYdP80b6GPWVlnDTkbfVWJqHl+2W8imHW+3482cUubMNAA5O5nEwGTDy4VwMMsJ0iqCLIizZGhw9i8n05ehb2A9b2dPvQD4W1um4lmJbeefEFBZNnb6XOFCLk0xgt5c7W+sSQAGtmEdFyz8qcp8Mvz4IoiD9VsdFqv4f5/Xtl7jbF2QIauiNqhFSOvrR351CjJaY6p3W+P8R/btYr/t/Y0Irl6/X4vTB3CN1g9ygnywHA0nkVE89td4QAOCHc0JuX++CLhFZuWxcPLeD/ch3+CjDlF6HAk6t54180nyQ3P0kWpii9f6MoFHvJRfBUTqAQclDdCZ0vLweNYwE6voA==")
		testPrivateKey := "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA0D1M1hg9Tjwn8mU/FtnKG9bGaomVWBFTeMeZIy3vX0HgNr1H\nNFCt7yHEE2H+BZkK8nJ2BHBeMw8hsmndmGuGyZrtWe8pxbl1x79/d+kcuG+cflJW\njBbx8Nkt8+SCImbU6BtW/jHoIvrKOhaAtxN8532A0OEzO5CWT8xoIE4yf9Xf5lle\nC2vb3zgjlQGaTzLu9O9Si9Zc1R//Oc16w1MoDo9nz8KHfzjjTR7K426zflMFbyb4\nEr+3Ywkj+BrzFvbtexltQaZfcfF9bWsdwXbhbM7QR3R9xLfJmGnohrHS6rQ/BXjW\nsQe9x73PjbP5bguYInn+DqTvUhiYw/KbyeOEQQIDAQABAoIBAQCnqrPTLnEuLQF9\nCkhh/bnd8HCSF3VIE6tB9HQ4/yNdb404he5vEQb7JBTcBmqh1zgZPlAIAvHV6rkX\nDmZ98xXz/epeH1NjAJD05BueUPPvDO7URzeoVFE5u6RkW/jr+iAzQtAom8ZtY8Cw\nRK4eunI3cbXmeWzm6OQeHFc6q7u9cONoo3bsALsHk3rEqEH1LUj4lEQPJPOD5g2H\ntvp1vs1SdjfM5nVcyvN4u4FBom97+vPDNtl12O6TBK2viVANjkw8FOVZIRqiAiNZ\nfIqOVwacN00phMzUiQB0IH0VRl+AVpKyWn5eAHTB6KPWGjF5LQxarJ8Km/DbYG3G\n7a1gEOiBAoGBAO9jZGMfRhpGeiOOM53ccUjzrTgNWIC8x+nGN69Z/mtzn5gg8fQO\nlypQcI8LVDPojqUkkBpxi7+6IuRXpzboC3yYmyvId6EBaD0ucCsJJbZfMuKoEkwP\nX9Hif90KkaJT6CXdvayzkLsoQVlFAUvIg6h1/YeT9UHDNxIjkKwfbN4vAoGBAN6w\nkNlQegU+UaCWs99uBm5mBqkl9TsP/fTwObkvnMF9RIaMDFWi7gZIDquihp8bicWs\n9qsc//BOJH+lJzN0ychklWFiy8GvRWH0nj80CkZIdRJBru6ssv616N/pF2UvS+2k\no8vK3Etr9Q9cu/TPdFk0DLdqGWwWK9gYQRSz8RiPAoGBAJOR4cB49u4bpA9nCcq2\nqd8e2BlFoNk7hsFFv+4IvB3hGPDe3khk9irPi5OimDWnlseW0n56oHuAcyHwJtRi\nFzKnoIBNA/HsvCV7Cwp8iRLzfJrcoOriT19DES9h5IT81I8DMnnT99Rn7GDrePEO\nmpquoauCOh5gCQLVicmRVbthAoGAVCHxF6lH8GMzA7DsFCXFWEBDk/Q7Si0ojTmV\nFVnfp1pkYVDX+CKuOsFOiZnFsqb8zioip1M1ftyG/ZKv1Mjy0zrtFPX2dR564B9D\nCi3nE9acJGGcbZ/hoEmpya6OoDPWQ9pH596ki/olg8BNYpheJLV9eG4lXKijt+ix\n7dht5hECgYB++7Y6z81/+FFB4MK4gyBvwOvPhk3J80Fx5KEx7kMBZ/i6VPUuWMKJ\n0r2YQGDf2K4DsECJ2eCDizkddnLSJvXUK8XuvHY+PsAGSLr6LeO8dxRdoPe+1ciT\nM4TC9j+kmeIJanPGR6wAKQ/CMblTAUQ76ztxscvqEAOUcfa2Y1n9xQ==\n-----END RSA PRIVATE KEY-----"
		testPublicKey := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0D1M1hg9Tjwn8mU/FtnK\nG9bGaomVWBFTeMeZIy3vX0HgNr1HNFCt7yHEE2H+BZkK8nJ2BHBeMw8hsmndmGuG\nyZrtWe8pxbl1x79/d+kcuG+cflJWjBbx8Nkt8+SCImbU6BtW/jHoIvrKOhaAtxN8\n532A0OEzO5CWT8xoIE4yf9Xf5lleC2vb3zgjlQGaTzLu9O9Si9Zc1R//Oc16w1Mo\nDo9nz8KHfzjjTR7K426zflMFbyb4Er+3Ywkj+BrzFvbtexltQaZfcfF9bWsdwXbh\nbM7QR3R9xLfJmGnohrHS6rQ/BXjWsQe9x73PjbP5bguYInn+DqTvUhiYw/KbyeOE\nQQIDAQAB\n-----END PUBLIC KEY-----"
		mockSecret := newMockSecret(t)
		testContext := context.Background()
		mockRegistry := newMockRegistry(t)
		mockDoguConfig := newMockConfigurationContext(t)
		mockDoguConfig.EXPECT().Get("public.pem").Return(testPublicKey, nil)
		//mockDoguConfig.On("Get", "public.pem").Return(testPublicKey)

		mockGlobalConfig := newMockGlobalConfigStore(t)
		mockGlobalConfig.EXPECT().Get("key_provider").Return(keyProviderStr, nil)
		mockRegistry.EXPECT().GlobalConfig().Return(mockGlobalConfig)
		mockRegistry.EXPECT().DoguConfig(doguname).Return(mockDoguConfig)
		encryptionAdapter := NewPublicKeyConfigEncryptionAdapter(mockSecret, mockRegistry, namespace)

		// when
		encryptedValue, err := encryptionAdapter.Encrypt(testContext, common.SimpleDoguName(doguname), common.SensitiveDoguConfigValue(testValue))

		// then
		require.NoError(t, err)

		// given
		mockSecret.EXPECT().Get(testContext, doguname+"-private", metav1.GetOptions{}).Return(&corev1.Secret{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Immutable:  nil,
			Data:       map[string][]byte{"private.pem": []byte(testPrivateKey)},
			StringData: nil,
			Type:       "",
		}, nil)
		// when
		decryptedValue, err := encryptionAdapter.Decrypt(testContext, common.SimpleDoguName(doguname), common.EncryptedDoguConfigValue(encryptedValue))
		// then
		require.NoError(t, err)
		assert.Equal(t, common.SensitiveDoguConfigValue(testValue), decryptedValue)
	})
}

func TestPublicKeyConfigEncryptionAdapter_Decrypt(t *testing.T) {
	t.Run("can decrypt", func(t *testing.T) {
		// given
		doguname := "testdogu"
		namespace := "testingnamespace"
		testValue := "value"
		keyProviderStr := "pkcs1v15"
		encryptedTestValue := common.EncryptedDoguConfigValue("Z5MYdP80b6GPWVlnDTkbfVWJqHl+2W8imHW+3482cUubMNAA5O5nEwGTDy4VwMMsJ0iqCLIizZGhw9i8n05ehb2A9b2dPvQD4W1um4lmJbeefEFBZNnb6XOFCLk0xgt5c7W+sSQAGtmEdFyz8qcp8Mvz4IoiD9VsdFqv4f5/Xtl7jbF2QIauiNqhFSOvrR351CjJaY6p3W+P8R/btYr/t/Y0Irl6/X4vTB3CN1g9ygnywHA0nkVE89td4QAOCHc0JuX++CLhFZuWxcPLeD/ch3+CjDlF6HAk6t54180nyQ3P0kWpii9f6MoFHvJRfBUTqAQclDdCZ0vLweNYwE6voA==")
		testPrivateKey := "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA0D1M1hg9Tjwn8mU/FtnKG9bGaomVWBFTeMeZIy3vX0HgNr1H\nNFCt7yHEE2H+BZkK8nJ2BHBeMw8hsmndmGuGyZrtWe8pxbl1x79/d+kcuG+cflJW\njBbx8Nkt8+SCImbU6BtW/jHoIvrKOhaAtxN8532A0OEzO5CWT8xoIE4yf9Xf5lle\nC2vb3zgjlQGaTzLu9O9Si9Zc1R//Oc16w1MoDo9nz8KHfzjjTR7K426zflMFbyb4\nEr+3Ywkj+BrzFvbtexltQaZfcfF9bWsdwXbhbM7QR3R9xLfJmGnohrHS6rQ/BXjW\nsQe9x73PjbP5bguYInn+DqTvUhiYw/KbyeOEQQIDAQABAoIBAQCnqrPTLnEuLQF9\nCkhh/bnd8HCSF3VIE6tB9HQ4/yNdb404he5vEQb7JBTcBmqh1zgZPlAIAvHV6rkX\nDmZ98xXz/epeH1NjAJD05BueUPPvDO7URzeoVFE5u6RkW/jr+iAzQtAom8ZtY8Cw\nRK4eunI3cbXmeWzm6OQeHFc6q7u9cONoo3bsALsHk3rEqEH1LUj4lEQPJPOD5g2H\ntvp1vs1SdjfM5nVcyvN4u4FBom97+vPDNtl12O6TBK2viVANjkw8FOVZIRqiAiNZ\nfIqOVwacN00phMzUiQB0IH0VRl+AVpKyWn5eAHTB6KPWGjF5LQxarJ8Km/DbYG3G\n7a1gEOiBAoGBAO9jZGMfRhpGeiOOM53ccUjzrTgNWIC8x+nGN69Z/mtzn5gg8fQO\nlypQcI8LVDPojqUkkBpxi7+6IuRXpzboC3yYmyvId6EBaD0ucCsJJbZfMuKoEkwP\nX9Hif90KkaJT6CXdvayzkLsoQVlFAUvIg6h1/YeT9UHDNxIjkKwfbN4vAoGBAN6w\nkNlQegU+UaCWs99uBm5mBqkl9TsP/fTwObkvnMF9RIaMDFWi7gZIDquihp8bicWs\n9qsc//BOJH+lJzN0ychklWFiy8GvRWH0nj80CkZIdRJBru6ssv616N/pF2UvS+2k\no8vK3Etr9Q9cu/TPdFk0DLdqGWwWK9gYQRSz8RiPAoGBAJOR4cB49u4bpA9nCcq2\nqd8e2BlFoNk7hsFFv+4IvB3hGPDe3khk9irPi5OimDWnlseW0n56oHuAcyHwJtRi\nFzKnoIBNA/HsvCV7Cwp8iRLzfJrcoOriT19DES9h5IT81I8DMnnT99Rn7GDrePEO\nmpquoauCOh5gCQLVicmRVbthAoGAVCHxF6lH8GMzA7DsFCXFWEBDk/Q7Si0ojTmV\nFVnfp1pkYVDX+CKuOsFOiZnFsqb8zioip1M1ftyG/ZKv1Mjy0zrtFPX2dR564B9D\nCi3nE9acJGGcbZ/hoEmpya6OoDPWQ9pH596ki/olg8BNYpheJLV9eG4lXKijt+ix\n7dht5hECgYB++7Y6z81/+FFB4MK4gyBvwOvPhk3J80Fx5KEx7kMBZ/i6VPUuWMKJ\n0r2YQGDf2K4DsECJ2eCDizkddnLSJvXUK8XuvHY+PsAGSLr6LeO8dxRdoPe+1ciT\nM4TC9j+kmeIJanPGR6wAKQ/CMblTAUQ76ztxscvqEAOUcfa2Y1n9xQ==\n-----END RSA PRIVATE KEY-----"
		mockSecret := newMockSecret(t)
		testContext := context.Background()
		mockRegistry := newMockRegistry(t)
		mockGlobalConfig := newMockGlobalConfigStore(t)
		mockGlobalConfig.EXPECT().Get("key_provider").Return(keyProviderStr, nil)
		mockRegistry.EXPECT().GlobalConfig().Return(mockGlobalConfig)
		//mockRegistry.EXPECT().DoguConfig(doguname).Return(mockDoguConfig)
		encryptionAdapter := NewPublicKeyConfigEncryptionAdapter(mockSecret, mockRegistry, namespace)
		mockSecret.EXPECT().Get(testContext, doguname+"-private", metav1.GetOptions{}).Return(&corev1.Secret{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Immutable:  nil,
			Data:       map[string][]byte{"private.pem": []byte(testPrivateKey)},
			StringData: nil,
			Type:       "",
		}, nil)

		// when
		decryptedValue, err := encryptionAdapter.Decrypt(testContext, common.SimpleDoguName(doguname), encryptedTestValue)

		// then
		require.NoError(t, err)
		assert.Equal(t, common.SensitiveDoguConfigValue(testValue), decryptedValue)
	})
}
