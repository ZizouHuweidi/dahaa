import * as React from 'react';
import { View } from 'react-native';
import { Link, useRouter } from 'expo-router';
import { Button } from '~/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '~/components/ui/card';
import { Input } from '~/components/ui/input';
import { Text } from '~/components/ui/text';
import { useAuthStore } from '~/lib/store/auth';

export default function RegisterScreen() {
  const router = useRouter();
  const { register, isLoading, error, clearError } = useAuthStore();
  const [username, setUsername] = React.useState('');
  const [email, setEmail] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [displayName, setDisplayName] = React.useState('');

  const handleRegister = async () => {
    if (!username.trim() || !email.trim() || !password.trim() || !displayName.trim()) return;
    await register(username, email, password, displayName);
    router.replace('/');
  };

  React.useEffect(() => {
    return () => clearError();
  }, []);

  return (
    <View className="flex-1 p-6 bg-background">
      <Card className="w-full max-w-sm mx-auto">
        <CardHeader>
          <CardTitle>Register</CardTitle>
          <CardDescription>
            Create a new account to start playing
          </CardDescription>
        </CardHeader>
        <CardContent className="gap-4">
          {error && (
            <Text className="text-destructive text-center">{error}</Text>
          )}
          <Input
            placeholder="Username"
            value={username}
            onChangeText={setUsername}
            autoCapitalize="none"
            autoComplete="username"
          />
          <Input
            placeholder="Email"
            value={email}
            onChangeText={setEmail}
            autoCapitalize="none"
            autoComplete="email"
            keyboardType="email-address"
          />
          <Input
            placeholder="Display Name"
            value={displayName}
            onChangeText={setDisplayName}
            autoComplete="name"
          />
          <Input
            placeholder="Password"
            value={password}
            onChangeText={setPassword}
            secureTextEntry
            autoCapitalize="none"
            autoComplete="password-new"
          />
        </CardContent>
        <CardFooter className="flex-col gap-4">
          <Button 
            className="w-full"
            onPress={handleRegister}
            disabled={isLoading || !username.trim() || !email.trim() || !password.trim() || !displayName.trim()}
          >
            <Text className="text-primary-foreground">
              {isLoading ? 'Creating account...' : 'Register'}
            </Text>
          </Button>
          <View className="flex-row justify-center">
            <Text className="text-muted-foreground">Already have an account? </Text>
            <Link href="/auth/login" asChild>
              <Text className="text-primary">Login</Text>
            </Link>
          </View>
        </CardFooter>
      </Card>
    </View>
  );
} 