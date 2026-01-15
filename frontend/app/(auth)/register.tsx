import { View, Text, TextInput, TouchableOpacity, StyleSheet, ScrollView, Alert } from 'react-native';
import { useState } from 'react';
import { router } from 'expo-router';

export default function Register() {
  const [form, setForm] = useState({
    name: '',
    username: '',
    email: '',
    password: '',
    gender: '',
    phone: '',
    country: '',
    referral_code: '',
  });

  const handleChange = (key: string, value: string) => {
    setForm({ ...form, [key]: value });
  };

  const [errors, setErrors] = useState<Record<string, string>>({});
  const validate = () => {
    const newErrors: Record<string, string> = {};

    if (!form.name) newErrors.name = 'Name has not been filled in!';
    if (!form.username) newErrors.username = 'Username has not been filled in!';
    if (!form.email) newErrors.email = 'Email has not been filled in!';
    if (!form.password) newErrors.password = 'Password has not been filled in!';
    if (!form.gender) newErrors.gender = 'Gender has not been filled in!';
    if (!form.phone) newErrors.phone = 'Phone has not been filled in!';
    if (!form.country) newErrors.country = 'Country has not been filled in!';

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleRegister = async () => {
    if (!validate()) return;

    try {
      const response = await fetch('http://192.168.1.35:8080/register/request-otp', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: form.name,
          username: form.username,
          email: form.email,
          password: form.password,
          gender: form.gender,
          phone: form.phone,
          country: form.country,
          referral_code: form.referral_code,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        Alert.alert('Register gagal', data.error || 'Terjadi kesalahan');
        return;
      }

      Alert.alert('Berhasil', 'OTP telah dikirim ke email');

      // PINDAH KE PAGE OTP
      router.push({
        pathname: '../(auth)/otp',
        params: { email: data.email },
      });

    } catch (error) {
      console.log(error);
      Alert.alert('Error', 'Tidak bisa terhubung ke server');
    }
  };

  return (
    <ScrollView contentContainerStyle={styles.container}>
      <View style={styles.card}>
        <Text style={styles.title}>Create Account</Text>

        {errors.name && <Text style={styles.errorText}>{errors.name}</Text>}
        <TextInput style={[styles.input, errors.name && styles.inputError,]} placeholder="Nama" onChangeText={(v) => {handleChange('name', v); setErrors({ ...errors, name: '' });}} />
        
        {errors.username && <Text style={styles.errorText}>{errors.username}</Text>}
        <TextInput style={[styles.input, errors.username && styles.inputError,]} placeholder="Username" onChangeText={(v) => {handleChange('username', v); setErrors({ ...errors, username: '' });}} />
        
        {errors.email && <Text style={styles.errorText}>{errors.email}</Text>}
        <TextInput style={[styles.input, errors.email && styles.inputError,]} placeholder="Email" autoCapitalize="none" onChangeText={(v) => {handleChange('email', v); setErrors({ ...errors, email: '' });}} />
        
        {errors.password && <Text style={styles.errorText}>{errors.password}</Text>}
        <TextInput style={[styles.input, errors.password && styles.inputError,]} placeholder="Password" secureTextEntry onChangeText={(v) => {handleChange('password', v); setErrors({ ...errors, password: '' });}} />
        
        {errors.gender && <Text style={styles.errorText}>{errors.gender}</Text>}
        <TextInput style={[styles.input, errors.gender && styles.inputError,]} placeholder="Gender" onChangeText={(v) => {handleChange('gender', v); setErrors({ ...errors, gender: '' });}} />
        
        {errors.phone && <Text style={styles.errorText}>{errors.phone}</Text>}
        <TextInput style={[styles.input, errors.phone && styles.inputError,]} placeholder="Phone" onChangeText={(v) => {handleChange('phone', v); setErrors({ ...errors, phone: '' });}} />
        
        {errors.country && <Text style={styles.errorText}>{errors.country}</Text>}
        <TextInput style={[styles.input, errors.country && styles.inputError,]} placeholder="Country" onChangeText={(v) => {handleChange('country', v); setErrors({ ...errors, country: '' });}} />
        
        <TextInput style={styles.input} placeholder="Referral Code (optional)" onChangeText={v => handleChange('referral_code', v)} />

        <TouchableOpacity style={styles.button} onPress={handleRegister}>
          <Text style={styles.buttonText}>Register</Text>
        </TouchableOpacity>

        <TouchableOpacity onPress={() => router.replace('/(auth)/login')}>
          <Text style={styles.link}>Sudah punya akun? Login</Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flexGrow: 1,
    backgroundColor: '#F5F7FA',
    justifyContent: 'center',
    padding: 20,
    paddingTop: 50,
  },
  card: {
    backgroundColor: '#FFF',
    borderRadius: 16,
    padding: 24,
    elevation: 4,
  },
  title: {
    fontSize: 22,
    fontWeight: '700',
    marginBottom: 20,
  },
  input: {
    height: 48,
    borderWidth: 1,
    borderColor: '#DDD',
    borderRadius: 10,
    paddingHorizontal: 14,
    marginBottom: 12,
  },
  button: {
    height: 48,
    backgroundColor: '#2563EB',
    borderRadius: 10,
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: 10,
  },
  buttonText: {
    color: '#FFF',
    fontSize: 16,
    fontWeight: '600',
  },
  link: {
    marginTop: 14,
    color: '#2563EB',
    textAlign: 'center',
  },
  inputError: {
    borderColor: '#EF4444',
  },

  errorText: {
    color: '#EF4444',
    fontSize: 12,
    marginBottom: 4,
    marginLeft: 4,
  },
});