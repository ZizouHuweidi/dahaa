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

export default function JoinGameScreen() {
  const router = useRouter();
  const { joinGame, isLoading, error, clearError } = useGameStore();
  const [playerName, setPlayerName] = React.useState('');
  const [gameCode, setGameCode] = React.useState('');

  const handleJoinGame = async () => {
    if (!playerName.trim() || !gameCode.trim()) return;
    await joinGame(gameCode, playerName);
    router.push('/game/lobby');
  };

  React.useEffect(() => {
    return () => clearError();
  }, []);

  return (
    <View className="flex-1 p-6 bg-background">
      <Card className="w-full max-w-sm mx-auto">
        <CardHeader>
          <CardTitle>Join Game</CardTitle>
          <CardDescription>
            Enter your name and the game code to join
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
          <Input
            placeholder="Game Code"
            value={gameCode}
            onChangeText={setGameCode}
            autoCapitalize="characters"
          />
        </CardContent>
        <CardFooter>
          <Button 
            className="w-full"
            onPress={handleJoinGame}
            disabled={isLoading || !playerName.trim() || !gameCode.trim()}
          >
            <Text className="text-primary-foreground">
              {isLoading ? 'Joining Game...' : 'Join Game'}
            </Text>
          </Button>
        </CardFooter>
      </Card>
    </View>
  );
} 