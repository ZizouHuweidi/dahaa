import * as React from 'react';
import { View } from 'react-native';
import { useRouter } from 'expo-router';
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
import { useGameStore } from '~/lib/store/game';

export default function CreateGameScreen() {
  const router = useRouter();
  const { createGame, isLoading, error, clearError } = useGameStore();
  const [playerName, setPlayerName] = React.useState('');

  const handleCreateGame = async () => {
    if (!playerName.trim()) return;
    await createGame(playerName);
    router.push('/game/lobby');
  };

  React.useEffect(() => {
    return () => clearError();
  }, []);

  return (
    <View className="flex-1 p-6 bg-background">
      <Card className="w-full max-w-sm mx-auto">
        <CardHeader>
          <CardTitle>Create New Game</CardTitle>
          <CardDescription>
            Enter your name to start a new game
          </CardDescription>
        </CardHeader>
        <CardContent className="gap-4">
          {error && (
            <Text className="text-destructive text-center">{error}</Text>
          )}
          <Input
            placeholder="Your Name"
            value={playerName}
            onChangeText={setPlayerName}
            autoComplete="name"
          />
        </CardContent>
        <CardFooter>
          <Button 
            className="w-full"
            onPress={handleCreateGame}
            disabled={isLoading || !playerName.trim()}
          >
            <Text className="text-primary-foreground">
              {isLoading ? 'Creating Game...' : 'Create Game'}
            </Text>
          </Button>
        </CardFooter>
      </Card>
    </View>
  );
} 