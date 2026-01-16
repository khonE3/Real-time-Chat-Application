// Isan Nong Bua Lam Phu Theme Colors
// Inspired by traditional Isan culture, temples, silk, and nature

export const isanColors = {
  // Primary - Temple Gold (สีทองวัด)
  // Inspired by the golden temples and Buddhist architecture
  gold: {
    50: '#FFF9E6',
    100: '#FFF0B3',
    200: '#FFE680',
    300: '#FFDB4D',
    400: '#FFD11A',
    500: '#D4A12A', // Primary
    600: '#B8860B',
    700: '#8B6914',
    800: '#5E4A0F',
    900: '#312709',
  },

  // Secondary - Silk Crimson (สีแดงผ้าไหม)
  // Inspired by traditional Isan silk weaving
  silk: {
    50: '#FFF0F0',
    100: '#FFD6D6',
    200: '#FFB3B3',
    300: '#FF8080',
    400: '#FF4D4D',
    500: '#DC143C', // Primary
    600: '#B22222',
    700: '#8B1A1A',
    800: '#661414',
    900: '#400D0D',
  },

  // Accent - Rice Paddy Green (สีเขียวทุ่งนา)
  // Inspired by the lush rice paddies of Isan
  paddy: {
    50: '#F0FFF0',
    100: '#C6F6C6',
    200: '#90EE90',
    300: '#5CDB5C',
    400: '#32CD32',
    500: '#228B22', // Primary
    600: '#1E7B1E',
    700: '#196619',
    800: '#145214',
    900: '#0D3D0D',
  },

  // Neutral - Earth Brown (สีน้ำตาลดิน)
  // Inspired by the laterite soil of the region
  earth: {
    50: '#FAF5F0',
    100: '#F0E6D6',
    200: '#E0CDB3',
    300: '#C9A87C',
    400: '#B8864B',
    500: '#8B4513', // Primary
    600: '#723A10',
    700: '#5A2E0D',
    800: '#42220A',
    900: '#2A1606',
  },

  // Background colors
  cotton: '#FFF8DC', // Handwoven cotton cream
  paper: '#FFFEF5', // Rice paper white
  
  // Text colors
  text: {
    primary: '#2D1B0E',
    secondary: '#5A4632',
    muted: '#8B7355',
    inverse: '#FFFEF5',
  },
};

// Tailwind CSS custom color configuration
export const tailwindColors = {
  gold: isanColors.gold,
  silk: isanColors.silk,
  paddy: isanColors.paddy,
  earth: isanColors.earth,
  cotton: isanColors.cotton,
  paper: isanColors.paper,
};

export default isanColors;
