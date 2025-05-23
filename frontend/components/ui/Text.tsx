import { Text as RNText, TextProps as RNTextProps, StyleSheet } from 'react-native';
import { useColorScheme } from 'react-native';

type TextProps = RNTextProps & {
  variant?: 'h1' | 'h2' | 'h3' | 'h4' | 'body' | 'caption' | 'button';
  weight?: 'regular' | 'medium' | 'semiBold' | 'bold';
  color?: 'primary' | 'secondary' | 'light' | 'error' | 'success';
  align?: 'left' | 'center' | 'right';
  children: React.ReactNode;
};

export function Text({
  variant = 'body',
  weight = 'regular',
  color = 'primary',
  align = 'left',
  style,
  children,
  ...props
}: TextProps) {
  const colorScheme = useColorScheme();
  const isDark = colorScheme === 'dark';

  const fontFamily = (() => {
    if (weight === 'medium') return 'Inter-Medium';
    if (weight === 'semiBold') return 'Inter-SemiBold';
    if (weight === 'bold') return 'Inter-Bold';
    return 'Inter-Regular';
  })();

  const fontSize = (() => {
    if (variant === 'h1') return 32;
    if (variant === 'h2') return 28;
    if (variant === 'h3') return 24;
    if (variant === 'h4') return 20;
    if (variant === 'caption') return 14;
    if (variant === 'button') return 16;
    return 16; // body
  })();

  const lineHeight = (() => {
    if (variant === 'h1') return 40;
    if (variant === 'h2') return 36;
    if (variant === 'h3') return 32;
    if (variant === 'h4') return 28;
    if (variant === 'caption') return 20;
    if (variant === 'button') return 24;
    return 24; // body
  })();

  const textColor = (() => {
    if (color === 'primary') return isDark ? '#F9FAFB' : '#111827';
    if (color === 'secondary') return isDark ? '#9CA3AF' : '#6B7280';
    if (color === 'light') return '#F9FAFB';
    if (color === 'error') return '#EF4444';
    if (color === 'success') return '#10B981';
    return isDark ? '#F9FAFB' : '#111827';
  })();

  return (
    <RNText
      style={[
        {
          fontFamily,
          fontSize,
          lineHeight,
          color: textColor,
          textAlign: align,
        },
        style,
      ]}
      {...props}
    >
      {children}
    </RNText>
  );
}
