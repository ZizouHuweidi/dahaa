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
import { Text } from '~/components/ui/text';
import { useGameStore } from '~/lib/store/game';

export default function ResultsScreen() {
  const router = useRouter();
  const { currentGame, isLoading, error } = useGameStore();

  // Determine the winner based on scores
  const winner = currentGame?.players?.reduce((prev, current) => {
    return (prev.score > current.score) ? prev : current;
  });

  if (!currentGame) {
    return (
      <View className="flex-1 p-6 bg-background justify-center items-center">
        <Text className="text-lg">No game results found.</Text>
        <Button className="mt-4" onPress={() => router.replace('/')}>Return Home</Button>
      </View>
    );
  }

  return (
    <View className="flex-1 p-6 bg-background">
      <Card className="w-full max-w-md mx-auto">
        <CardHeader>
          <CardTitle>Game Results</CardTitle>
          <CardDescription>Final scores and winner</CardDescription>
        </CardHeader>
        <CardContent className="gap-4">
          {error && <Text className="text-destructive text-center">{error}</Text>}
          <View className="items-center mb-4">
            <Text className="text-xl font-bold">Winner:</Text>
            <Text className="text-2xl font-bold text-primary">{winner?.name}</Text>
          </View>
          <View>
            <Text className="font-semibold mb-2">Final Scores:</Text>
            {currentGame.players.map((player) => (
              <Text key={player.id} className="text-base">
                {player.name}: {player.score} points
              </Text>
            ))}
          </View>
        </CardContent>
        <CardFooter>
          <Button className="w-full" onPress={() => router.replace('/')}>
            <Text className="text-primary-foreground">Return Home</Text>
          </Button>
        </CardFooter>
      </Card>
    </View>
  );
} 