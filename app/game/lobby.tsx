import * as React from 'react';
import { View, ScrollView } from 'react-native';
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
import { Text } from '~/components/ui/text';
import { useGameStore, Player } from '~/lib/store/game';
import { useAuthStore } from '~/lib/store/auth';

export default function GameLobbyScreen() {
  const router = useRouter();
  const { currentGame, isLoading, error, startGame } = useGameStore();
  const { user } = useAuthStore();
  const [isStarting, setIsStarting] = React.useState(false);

  const handleStartGame = async () => {
    if (!currentGame) return;
    setIsStarting(true);
    try {
      await startGame();
      router.push('/game/play');
    } catch (err) {
      console.error('Failed to start game:', err);
    } finally {
      setIsStarting(false);
    }
  };

  const handleLeaveGame = () => {
    router.replace('/');
  };

  if (!currentGame) {
    return (
      <View className="flex-1 p-6 bg-background">
        <Card className="w-full max-w-sm mx-auto">
          <CardHeader>
            <CardTitle>Game Not Found</CardTitle>
            <CardDescription>
              The game you're looking for doesn't exist or you've been disconnected.
            </CardDescription>
          </CardHeader>
          <CardFooter>
            <Button className="w-full" onPress={() => router.replace('/')}>
              <Text className="text-primary-foreground">Return Home</Text>
            </Button>
          </CardFooter>
        </Card>
      </View>
    );
  }

  return (
    <View className="flex-1 p-6 bg-background">
      <Card className="w-full max-w-sm mx-auto">
        <CardHeader>
          <CardTitle>Game Lobby</CardTitle>
          <CardDescription>
            Game Code: {currentGame.code}
          </CardDescription>
        </CardHeader>
        <CardContent className="gap-4">
          {error && (
            <Text className="text-destructive text-center">{error}</Text>
          )}
          <View>
            <Text className="text-lg font-semibold mb-2">Players:</Text>
            <ScrollView className="max-h-48">
              {currentGame.players.map((player: Player) => (
                <View
                  key={player.id}
                  className="flex-row items-center justify-between py-2 border-b border-border"
                >
                  <Text>{player.name}</Text>
                </View>
              ))}
            </ScrollView>
          </View>
        </CardContent>
        <CardFooter className="flex-col gap-2">
          {user && currentGame.hostId === user.id && (
            <Button
              className="w-full"
              onPress={handleStartGame}
              disabled={isStarting || currentGame.players.length < 2}
            >
              <Text className="text-primary-foreground">
                {isStarting ? 'Starting Game...' : 'Start Game'}
              </Text>
            </Button>
          )}
          <Button
            className="w-full"
            variant="outline"
            onPress={handleLeaveGame}
            disabled={isLoading}
          >
            <Text>Leave Game</Text>
          </Button>
        </CardFooter>
      </Card>
    </View>
  );
} 