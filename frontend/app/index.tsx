import { Redirect } from 'expo-router';
import { useAuthStore } from '../src/authStore/authStore';

export default function Index() {
  const token = useAuthStore((s) => s.token);

  if (!token) {
    return <Redirect href="../(auth)/register" />;
  }

  return <Redirect href="../(tabs)" />;
}