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

export default function LoginScreen() {
  const router = useRouter();
  const { login, isLoading, error, clearError } = useAuthStore();
  const [username, setUsername] = React.useState('');
  const [password, setPassword] = React.useState('');

  const handleLogin = async () => {
    if (!username.trim() || !password.trim()) return;
    await login(username, password);
    router.replace('/');
  };

  React.useEffect(() => {
    return () => clearError();
  }, []);

  return (
    <View className="flex-1 p-6 bg-background">
      <Card className="w-full max-w-sm mx-auto">
        <CardHeader>
          <CardTitle>Login</CardTitle>
          <CardDescription>
            Enter your credentials to continue
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
            placeholder="Password"
            value={password}
            onChangeText={setPassword}
            secureTextEntry
            autoCapitalize="none"
            autoComplete="password"
          />
        </CardContent>
        <CardFooter className="flex-col gap-4">
          <Button 
            className="w-full"
            onPress={handleLogin}
            disabled={isLoading || !username.trim() || !password.trim()}
          >
            <Text className="text-primary-foreground">
              {isLoading ? 'Logging in...' : 'Login'}
            </Text>
          </Button>
          <View className="flex-row justify-center">
            <Text className="text-muted-foreground">Don't have an account? </Text>
            <Link href="/auth/register" asChild>
              <Text className="text-primary">Register</Text>
            </Link>
          </View>
        </CardFooter>
      </Card>
    </View>
  );
} 