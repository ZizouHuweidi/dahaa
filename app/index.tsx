import * as React from 'react';
import { View } from 'react-native';
import { Link } from 'expo-router';
import { GameIcon } from '~/lib/icons/Game';
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

export default function HomeScreen() {
  return (
    <View className="flex-1 justify-center items-center gap-5 p-6 bg-background">
      <Card className="w-full max-w-sm p-6 rounded-2xl">
        <CardHeader className="items-center">
          <View className="w-24 h-24 items-center justify-center rounded-full bg-primary/10">
            <GameIcon size={48} className="text-primary" />
          </View>
          <View className="p-3" />
          <CardTitle className="pb-2 text-center text-2xl">Dahaa</CardTitle>
          <CardDescription className="text-center text-base">
            The ultimate trivia game where you create the answers
          </CardDescription>
        </CardHeader>
        <CardContent className="gap-4">
          <Link href="/game/create" asChild>
            <Button className="w-full">
              <Text className="text-primary-foreground">Create Game</Text>
            </Button>
          </Link>
          <Link href="/game/join" asChild>
            <Button variant="outline" className="w-full">
              <Text>Join Game</Text>
            </Button>
          </Link>
        </CardContent>
        <CardFooter className="flex-col gap-2">
          <Text className="text-sm text-muted-foreground text-center">
            Create a new game or join an existing one using a game code
          </Text>
        </CardFooter>
      </Card>
    </View>
  );
}
