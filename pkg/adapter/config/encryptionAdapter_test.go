package config

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var testPrivateKey1 = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA0D1M1hg9Tjwn8mU/FtnKG9bGaomVWBFTeMeZIy3vX0HgNr1H\nNFCt7yHEE2H+BZkK8nJ2BHBeMw8hsmndmGuGyZrtWe8pxbl1x79/d+kcuG+cflJW\njBbx8Nkt8+SCImbU6BtW/jHoIvrKOhaAtxN8532A0OEzO5CWT8xoIE4yf9Xf5lle\nC2vb3zgjlQGaTzLu9O9Si9Zc1R//Oc16w1MoDo9nz8KHfzjjTR7K426zflMFbyb4\nEr+3Ywkj+BrzFvbtexltQaZfcfF9bWsdwXbhbM7QR3R9xLfJmGnohrHS6rQ/BXjW\nsQe9x73PjbP5bguYInn+DqTvUhiYw/KbyeOEQQIDAQABAoIBAQCnqrPTLnEuLQF9\nCkhh/bnd8HCSF3VIE6tB9HQ4/yNdb404he5vEQb7JBTcBmqh1zgZPlAIAvHV6rkX\nDmZ98xXz/epeH1NjAJD05BueUPPvDO7URzeoVFE5u6RkW/jr+iAzQtAom8ZtY8Cw\nRK4eunI3cbXmeWzm6OQeHFc6q7u9cONoo3bsALsHk3rEqEH1LUj4lEQPJPOD5g2H\ntvp1vs1SdjfM5nVcyvN4u4FBom97+vPDNtl12O6TBK2viVANjkw8FOVZIRqiAiNZ\nfIqOVwacN00phMzUiQB0IH0VRl+AVpKyWn5eAHTB6KPWGjF5LQxarJ8Km/DbYG3G\n7a1gEOiBAoGBAO9jZGMfRhpGeiOOM53ccUjzrTgNWIC8x+nGN69Z/mtzn5gg8fQO\nlypQcI8LVDPojqUkkBpxi7+6IuRXpzboC3yYmyvId6EBaD0ucCsJJbZfMuKoEkwP\nX9Hif90KkaJT6CXdvayzkLsoQVlFAUvIg6h1/YeT9UHDNxIjkKwfbN4vAoGBAN6w\nkNlQegU+UaCWs99uBm5mBqkl9TsP/fTwObkvnMF9RIaMDFWi7gZIDquihp8bicWs\n9qsc//BOJH+lJzN0ychklWFiy8GvRWH0nj80CkZIdRJBru6ssv616N/pF2UvS+2k\no8vK3Etr9Q9cu/TPdFk0DLdqGWwWK9gYQRSz8RiPAoGBAJOR4cB49u4bpA9nCcq2\nqd8e2BlFoNk7hsFFv+4IvB3hGPDe3khk9irPi5OimDWnlseW0n56oHuAcyHwJtRi\nFzKnoIBNA/HsvCV7Cwp8iRLzfJrcoOriT19DES9h5IT81I8DMnnT99Rn7GDrePEO\nmpquoauCOh5gCQLVicmRVbthAoGAVCHxF6lH8GMzA7DsFCXFWEBDk/Q7Si0ojTmV\nFVnfp1pkYVDX+CKuOsFOiZnFsqb8zioip1M1ftyG/ZKv1Mjy0zrtFPX2dR564B9D\nCi3nE9acJGGcbZ/hoEmpya6OoDPWQ9pH596ki/olg8BNYpheJLV9eG4lXKijt+ix\n7dht5hECgYB++7Y6z81/+FFB4MK4gyBvwOvPhk3J80Fx5KEx7kMBZ/i6VPUuWMKJ\n0r2YQGDf2K4DsECJ2eCDizkddnLSJvXUK8XuvHY+PsAGSLr6LeO8dxRdoPe+1ciT\nM4TC9j+kmeIJanPGR6wAKQ/CMblTAUQ76ztxscvqEAOUcfa2Y1n9xQ==\n-----END RSA PRIVATE KEY-----"
var testPublicKey1 = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0D1M1hg9Tjwn8mU/FtnK\nG9bGaomVWBFTeMeZIy3vX0HgNr1HNFCt7yHEE2H+BZkK8nJ2BHBeMw8hsmndmGuG\nyZrtWe8pxbl1x79/d+kcuG+cflJWjBbx8Nkt8+SCImbU6BtW/jHoIvrKOhaAtxN8\n532A0OEzO5CWT8xoIE4yf9Xf5lleC2vb3zgjlQGaTzLu9O9Si9Zc1R//Oc16w1Mo\nDo9nz8KHfzjjTR7K426zflMFbyb4Er+3Ywkj+BrzFvbtexltQaZfcfF9bWsdwXbh\nbM7QR3R9xLfJmGnohrHS6rQ/BXjWsQe9x73PjbP5bguYInn+DqTvUhiYw/KbyeOE\nQQIDAQAB\n-----END PUBLIC KEY-----"
var testPrivateKey2 = "-----BEGIN RSA PRIVATE KEY-----\nMIIEoAIBAAKCAQEA5J8NgPLcg6xN3HXH+fkbDrl2UnPEmxb+GrlCzXKT3475uqoC\nzxo327I+8r4yhUkRGm1QzUr1msk67npsybuZP6h4+pnJI8aPc/xTDx6dXvUjKOxf\nWYQ5zYylG9lajFBPxkU0bBy4pfFu3titlthnLzpINxIwYXoGtIPNhnyRif308It+\n8N8r4rVrSkEc3owMS6o5zv53K2hxftdnYrPGWfEm39pYg+Xep+lFdTqX3Dpxx0rs\nsxXenAaYgxjsIi9zwt7yCvv/TyibF0ZYDIVsoqgxfrlSoLd2J1cxjStPiPfUNiyH\nLCDm8HwI88pQuz9UIlGgMUVIf3MxgqOD12dg/QIDAQABAoIBAFl6btyTMP9QBsFM\nT9JkTtS6fbbTnJVesGFhNOYX/Aw5d0A5nhPUnRwdbUmwazGDYXBIbKGMvwevzqLb\nw4xJIjeqBn9+hRy9cBPjI9b9EnbB1tsDeGYevEjYzR5TOX9FR5PALj5KF3LLRatu\nfrJVTD1NwEndkpX0Hn+0PlJumr+4qhoTpMb5hwr9lUKJ96fM2MU6vgG75J2d84FJ\n4AAhBsgIX7avjA+0efml/vYYdVk3H7/Vq6yvTst+pUgT40d11pcJUGNZJQlUJWAQ\nYFhwbXZQTbYOp3Av89Q2+bnPheFEnTDxQ6HGNS+fix/DQamEAne1pDUyMDvdr2Df\ntzXm0bkCgYEA6sWN2EbWLziIMaFFLL+7r4LRQERc75UYr5/xTKbph/L+IjqNQ4zN\n0rxlEKVUjhERZjdepOc2pKRSx+nXFW4myv3tV2CInwgb/vmPjsxDtwXyECo+84MH\nYNP+wnRkSxTZvIEMVjtwSv8WrrSnmvjTBTkPkdvY0AnA44gn50gqZu8CgYEA+Ush\nbr9KYQskGb3AVRr/EMWy1j7mXyap+ErBH68QiZCTkS49tFdmhDlI6n8P8SUesht5\n6KOAQdL/6x6gAE764ro0qLEaNNZv9/6Y5oF2Kv7Xlai/8NKuL9UGLNu5YxnPGi/m\nwfJUClQ8TiMgYVtJrRHv/0E4HlXRz4UUyqIkFtMCgYANrG3jf9SvsWI1SchGn/Al\ne8AGNzUWex+R8wXRyhLl6SAmDDT4DzZZpMFaI9b140aZJnZrsk+7bRqpLBRr2huG\nTR3KrgOnB4jh49UZown6me0MRfmeoy4F1LMMzkydFtzLntSCHTogFBVVHY55dy6L\nKlSe0Sgijb7fQanZTZmynwKBgGpulRt/N/Yul38V8CNlnzg975hgymIdU7vZzpIE\nX/8bZqU5JMb1+aLCAkt7bAb8XhqUeHvGMl/oAbMUJCN9lMdv0EOlORcN5kfuvsDK\nzPSWUNxoa2oZyJxLSpOkS4Xv4ue/Q7nSB+dRB14kyRJHszDc06Ya5iatZSJAIxxQ\nFTBZAn9MMYMxue8HYEjXnl+onsHkNrrnP+AZS9vKEggfOb2A1J8rDbYjVCyMz+BZ\nrCnLNeCZYf3aPHh2UJozntZr51ylTuzb01TI92DbVXY1Rga/efDDD6zJpvT51zQf\nGyMn91fjbl5aP7DHlPjQYQtf62aJRHHebOnIXg/RDyPUsmc1\n-----END RSA PRIVATE KEY-----"
var testPublicKey2 = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5J8NgPLcg6xN3HXH+fkb\nDrl2UnPEmxb+GrlCzXKT3475uqoCzxo327I+8r4yhUkRGm1QzUr1msk67npsybuZ\nP6h4+pnJI8aPc/xTDx6dXvUjKOxfWYQ5zYylG9lajFBPxkU0bBy4pfFu3titlthn\nLzpINxIwYXoGtIPNhnyRif308It+8N8r4rVrSkEc3owMS6o5zv53K2hxftdnYrPG\nWfEm39pYg+Xep+lFdTqX3Dpxx0rssxXenAaYgxjsIi9zwt7yCvv/TyibF0ZYDIVs\noqgxfrlSoLd2J1cxjStPiPfUNiyHLCDm8HwI88pQuz9UIlGgMUVIf3MxgqOD12dg\n/QIDAQAB\n-----END PUBLIC KEY-----"

func TestPublicKeyConfigEncryptionAdapter_Encrypt(t *testing.T) {
	t.Run("can encrypt and decrypt", func(t *testing.T) {
		// given
		doguname := "testdogu"
		namespace := "testingnamespace"
		testValue := "value"
		keyProviderStr := "pkcs1v15"
		mockSecret := newMockSecret(t)
		testContext := context.Background()
		mockRegistry := newMockRegistry(t)
		mockDoguConfig := newMockConfigurationContext(t)
		mockDoguConfig.EXPECT().Get("public.pem").Return(testPublicKey1, nil)
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
			Data:       map[string][]byte{"private.pem": []byte(testPrivateKey1)},
			StringData: nil,
			Type:       "",
		}, nil)
		// when
		decryptedValue, err := encryptionAdapter.Decrypt(testContext, common.SimpleDoguName(doguname), common.EncryptedDoguConfigValue(encryptedValue))
		// then
		require.NoError(t, err)
		assert.Equal(t, common.SensitiveDoguConfigValue(testValue), decryptedValue)
	})

	t.Run("can not get pub key", func(t *testing.T) {
		// given
		doguname := "testdogu"
		namespace := "testingnamespace"
		testValue := "value"
		mockSecret := newMockSecret(t)
		testContext := context.Background()
		mockRegistry := newMockRegistry(t)
		mockGlobalConfig := newMockGlobalConfigStore(t)
		mockGlobalConfig.EXPECT().Get("key_provider").Return("", fmt.Errorf("nope"))
		mockRegistry.EXPECT().GlobalConfig().Return(mockGlobalConfig)
		encryptionAdapter := NewPublicKeyConfigEncryptionAdapter(mockSecret, mockRegistry, namespace)

		// when
		encryptedValue, err := encryptionAdapter.Encrypt(testContext, common.SimpleDoguName(doguname), common.SensitiveDoguConfigValue(testValue))

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not get public key")
		assert.Equal(t, common.EncryptedDoguConfigValue(""), encryptedValue)
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
		mockSecret := newMockSecret(t)
		testContext := context.Background()
		mockRegistry := newMockRegistry(t)
		mockGlobalConfig := newMockGlobalConfigStore(t)
		mockGlobalConfig.EXPECT().Get("key_provider").Return(keyProviderStr, nil)
		mockRegistry.EXPECT().GlobalConfig().Return(mockGlobalConfig)
		encryptionAdapter := NewPublicKeyConfigEncryptionAdapter(mockSecret, mockRegistry, namespace)
		mockSecret.EXPECT().Get(testContext, doguname+"-private", metav1.GetOptions{}).Return(&corev1.Secret{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Immutable:  nil,
			Data:       map[string][]byte{"private.pem": []byte(testPrivateKey1)},
			StringData: nil,
			Type:       "",
		}, nil)

		// when
		decryptedValue, err := encryptionAdapter.Decrypt(testContext, common.SimpleDoguName(doguname), encryptedTestValue)

		// then
		require.NoError(t, err)
		assert.Equal(t, common.SensitiveDoguConfigValue(testValue), decryptedValue)
	})

	t.Run("can not get key pair", func(t *testing.T) {
		// given
		doguname := "testdogu"
		namespace := "testingnamespace"
		encryptedTestValue := common.EncryptedDoguConfigValue("Z5MYdP80b6GPWVlnDTkbfVWJqHl+2W8imHW+3482cUubMNAA5O5nEwGTDy4VwMMsJ0iqCLIizZGhw9i8n05ehb2A9b2dPvQD4W1um4lmJbeefEFBZNnb6XOFCLk0xgt5c7W+sSQAGtmEdFyz8qcp8Mvz4IoiD9VsdFqv4f5/Xtl7jbF2QIauiNqhFSOvrR351CjJaY6p3W+P8R/btYr/t/Y0Irl6/X4vTB3CN1g9ygnywHA0nkVE89td4QAOCHc0JuX++CLhFZuWxcPLeD/ch3+CjDlF6HAk6t54180nyQ3P0kWpii9f6MoFHvJRfBUTqAQclDdCZ0vLweNYwE6voA==")
		mockSecret := newMockSecret(t)
		testContext := context.Background()
		mockRegistry := newMockRegistry(t)
		mockGlobalConfig := newMockGlobalConfigStore(t)
		mockGlobalConfig.EXPECT().Get("key_provider").Return("", fmt.Errorf("nope"))
		mockRegistry.EXPECT().GlobalConfig().Return(mockGlobalConfig)
		encryptionAdapter := NewPublicKeyConfigEncryptionAdapter(mockSecret, mockRegistry, namespace)
		mockSecret.EXPECT().Get(testContext, doguname+"-private", metav1.GetOptions{}).Return(&corev1.Secret{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Immutable:  nil,
			Data:       map[string][]byte{"private.pem": []byte(testPrivateKey1)},
			StringData: nil,
			Type:       "",
		}, nil)

		// when
		decryptedValue, err := encryptionAdapter.Decrypt(testContext, common.SimpleDoguName(doguname), encryptedTestValue)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not get key pair for dogu "+doguname)
		assert.Equal(t, common.SensitiveDoguConfigValue(""), decryptedValue)
	})
}

func TestPublicKeyConfigEncryptionAdapter_EncryptAll(t *testing.T) {
	t.Run("can encrypt and decrypt multiple values", func(t *testing.T) {
		// given
		doguname1 := "testdogu1"
		doguname2 := "testdogu2"
		testkey1 := "testkey1"
		testkey2 := "testkey2"
		testvalue1 := "testvalue1"
		testvalue2 := "testvalue2"
		namespace := "testingnamespace"
		testValues := map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
			common.SensitiveDoguConfigKey{common.DoguConfigKey{
				DoguName: common.SimpleDoguName(doguname1),
				Key:      testkey1,
			}}: common.SensitiveDoguConfigValue(testvalue1),
			common.SensitiveDoguConfigKey{common.DoguConfigKey{
				DoguName: common.SimpleDoguName(doguname2),
				Key:      testkey2,
			}}: common.SensitiveDoguConfigValue(testvalue2),
		}
		keyProviderStr := "pkcs1v15"

		mockSecret := newMockSecret(t)
		testContext := context.Background()
		mockRegistry := newMockRegistry(t)
		mockDoguConfig1 := newMockConfigurationContext(t)
		mockDoguConfig2 := newMockConfigurationContext(t)
		mockDoguConfig1.EXPECT().Get("public.pem").Return(testPublicKey1, nil)
		mockDoguConfig2.EXPECT().Get("public.pem").Return(testPublicKey2, nil)
		mockGlobalConfig := newMockGlobalConfigStore(t)
		mockGlobalConfig.EXPECT().Get("key_provider").Return(keyProviderStr, nil)
		mockRegistry.EXPECT().GlobalConfig().Return(mockGlobalConfig)
		mockRegistry.EXPECT().DoguConfig(doguname1).Return(mockDoguConfig1)
		mockRegistry.EXPECT().DoguConfig(doguname2).Return(mockDoguConfig2)
		encryptionAdapter := NewPublicKeyConfigEncryptionAdapter(mockSecret, mockRegistry, namespace)

		// when
		encryptedValues, err := encryptionAdapter.EncryptAll(testContext, testValues)

		// then
		require.NoError(t, err)

		// given
		mockSecret.EXPECT().Get(testContext, doguname1+"-private", metav1.GetOptions{}).Return(&corev1.Secret{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Immutable:  nil,
			Data:       map[string][]byte{"private.pem": []byte(testPrivateKey1)},
			StringData: nil,
			Type:       "",
		}, nil)
		mockSecret.EXPECT().Get(testContext, doguname2+"-private", metav1.GetOptions{}).Return(&corev1.Secret{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Immutable:  nil,
			Data:       map[string][]byte{"private.pem": []byte(testPrivateKey2)},
			StringData: nil,
			Type:       "",
		}, nil)
		// when
		decryptedValues, err := encryptionAdapter.DecryptAll(testContext, encryptedValues)
		// then
		require.NoError(t, err)
		assert.Equal(t, common.SensitiveDoguConfigValue(testvalue1), decryptedValues[common.SensitiveDoguConfigKey{common.DoguConfigKey{DoguName: common.SimpleDoguName(doguname1), Key: testkey1}}])
		assert.Equal(t, common.SensitiveDoguConfigValue(testvalue2), decryptedValues[common.SensitiveDoguConfigKey{common.DoguConfigKey{DoguName: common.SimpleDoguName(doguname2), Key: testkey2}}])
	})
}
