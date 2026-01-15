import { View, Text } from 'react-native';
import { useLocalSearchParams } from 'expo-router';

export default function OTP() {
  const { email } = useLocalSearchParams();

  return (
    <View>
      <Text>Verifikasi OTP</Text>
      <Text>Email: {email}</Text>
    </View>
  );
}