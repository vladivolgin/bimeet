import { extendTheme } from '@chakra-ui/react'

const theme = extendTheme({
  config: {
    initialColorMode: 'light',
    useSystemColorMode: false,
  },
  fonts: {
    heading: `'Inter', system-ui, sans-serif`,
    body: `'Inter', system-ui, sans-serif`,
  },
  colors: {
    // Indigo — основной бренд-цвет (Tailwind Indigo scale)
    brand: {
      50:  '#EEF2FF',
      100: '#E0E7FF',
      200: '#C7D2FE',
      300: '#A5B4FC',
      400: '#818CF8',
      500: '#6366F1',
      600: '#4F46E5',
      700: '#4338CA',
      800: '#3730A3',
      900: '#312E81',
    },
    // Slate — для типографики
    slate: {
      50:  '#F8FAFC',
      100: '#F1F5F9',
      200: '#E2E8F0',
      300: '#CBD5E1',
      400: '#94A3B8',
      500: '#64748B',
      600: '#475569',
      700: '#334155',
      800: '#1E293B',
      900: '#0F172A',
    },
  },
  // Semantic tokens — автоматически переключаются между light/dark
  semanticTokens: {
    colors: {
      pageBg: {
        default: '#F8F7F4',
        _dark: '#0E0F16',
      },
      cardBg: {
        default: 'white',
        _dark: '#171923',
      },
      subtleBg: {
        default: '#F8FAFC',
        _dark: '#1D1F2E',
      },
      inputBg: {
        default: 'white',
        _dark: '#1A1C2A',
      },
      cardBorder: {
        default: 'rgba(15,23,42,0.07)',
        _dark: 'rgba(255,255,255,0.08)',
      },
      defaultBorder: {
        default: '#E2E8F0',
        _dark: 'rgba(255,255,255,0.12)',
      },
      subtleBorder: {
        default: '#F1F5F9',
        _dark: 'rgba(255,255,255,0.07)',
      },
      mainText: {
        default: '#0F172A',
        _dark: '#F1F5F9',
      },
      dimText: {
        default: '#64748B',
        _dark: '#94A3B8',
      },
      faintText: {
        default: '#94A3B8',
        _dark: '#64748B',
      },
      sidebarBg: {
        default: '#FFFFFF',
        _dark: '#13141F',
      },
      sidebarBorder: {
        default: 'rgba(15,23,42,0.08)',
        _dark: 'rgba(255,255,255,0.07)',
      },
      navActiveBg: {
        default: '#EEF2FF',
        _dark: 'rgba(99,102,241,0.15)',
      },
      navActiveText: {
        default: '#3730A3',
        _dark: '#A5B4FC',
      },
    },
  },
  styles: {
    global: {
      body: {
        bg: 'pageBg',
        color: 'mainText',
      },
    },
  },
  shadows: {
    card: '0 1px 3px rgba(15,23,42,0.04), 0 4px 16px rgba(15,23,42,0.05)',
    cardHover: '0 2px 8px rgba(15,23,42,0.06), 0 12px 32px rgba(15,23,42,0.09)',
    brand: '0 4px 16px rgba(79, 70, 229, 0.28)',
  },
  components: {
    Heading: {
      baseStyle: {
        fontWeight: '600',
        letterSpacing: '-0.02em',
        color: 'mainText',
      },
      sizes: {
        xl:  { fontSize: '1.5rem',    lineHeight: '1.3' },
        lg:  { fontSize: '1.1875rem', lineHeight: '1.35' },
        md:  { fontSize: '1rem',      lineHeight: '1.4' },
        sm:  { fontSize: '0.9375rem', lineHeight: '1.4' },
        xs:  { fontSize: '0.8125rem', lineHeight: '1.35' },
      },
    },
    Button: {
      baseStyle: {
        borderRadius: '10px',
        fontWeight: 500,
        fontSize: 'sm',
        transition: 'all 0.2s',
        letterSpacing: '0.005em',
      },
      variants: {
        solid: (props: { colorScheme: string }) => {
          if (props.colorScheme === 'brand') {
            return {
              bg: 'brand.600',
              color: 'white',
              _hover: {
                bg: 'brand.700',
                transform: 'translateY(-1px)',
                boxShadow: '0 4px 14px rgba(79, 70, 229, 0.35)',
              },
              _active: {
                transform: 'translateY(0)',
                bg: 'brand.800',
                boxShadow: 'none',
              },
            }
          }
          if (props.colorScheme === 'blue') {
            return {
              bg: 'brand.600',
              color: 'white',
              _hover: {
                bg: 'brand.700',
                transform: 'translateY(-1px)',
                boxShadow: '0 4px 14px rgba(79, 70, 229, 0.35)',
              },
              _active: {
                transform: 'translateY(0)',
                bg: 'brand.800',
                boxShadow: 'none',
              },
            }
          }
          return {}
        },
        ghost: {
          borderRadius: '10px',
          color: 'dimText',
          _hover: { bg: 'rgba(15,23,42,0.05)', color: 'mainText' },
          _dark: {
            _hover: { bg: 'rgba(255,255,255,0.07)', color: 'mainText' },
          },
        },
        outline: {
          borderRadius: '10px',
          borderWidth: '1px',
          borderColor: 'defaultBorder',
          color: '#334155',
          _hover: { bg: '#F8FAFC', borderColor: '#A5B4FC' },
          _dark: {
            color: '#CBD5E1',
            _hover: { bg: 'rgba(255,255,255,0.05)', borderColor: 'brand.400' },
          },
        },
      },
    },
    Input: {
      variants: {
        outline: {
          field: {
            borderRadius: '10px',
            bg: 'inputBg',
            borderColor: 'defaultBorder',
            fontSize: 'sm',
            color: 'mainText',
            _placeholder: { color: 'faintText' },
            _focus: {
              borderColor: 'brand.500',
              boxShadow: '0 0 0 3px rgba(79, 70, 229, 0.12)',
            },
            _hover: { borderColor: '#A5B4FC' },
            _dark: {
              _focus: {
                boxShadow: '0 0 0 3px rgba(99, 102, 241, 0.20)',
              },
            },
          },
        },
      },
      defaultProps: { variant: 'outline' },
    },
    Textarea: {
      variants: {
        outline: {
          borderRadius: '10px',
          bg: 'inputBg',
          borderColor: 'defaultBorder',
          fontSize: 'sm',
          color: 'mainText',
          _placeholder: { color: 'faintText' },
          _focus: {
            borderColor: 'brand.500',
            boxShadow: '0 0 0 3px rgba(79, 70, 229, 0.12)',
          },
          _hover: { borderColor: '#A5B4FC' },
          _dark: {
            _focus: {
              boxShadow: '0 0 0 3px rgba(99, 102, 241, 0.20)',
            },
          },
        },
      },
    },
    NumberInput: {
      variants: {
        outline: {
          field: {
            borderRadius: '10px',
            bg: 'inputBg',
            fontSize: 'sm',
          },
        },
      },
    },
    Card: {
      baseStyle: {
        container: {
          borderRadius: 'xl',
          boxShadow: 'card',
          border: '1px solid',
          borderColor: 'cardBorder',
          bg: 'cardBg',
          overflow: 'hidden',
        },
      },
    },
    Badge: {
      baseStyle: {
        borderRadius: 'md',
        px: 2.5,
        py: 0.5,
        fontWeight: 500,
        fontSize: 'xs',
        textTransform: 'none',
        letterSpacing: '0',
      },
    },
    Tabs: {
      variants: {
        line: {
          tab: {
            fontWeight: 450,
            fontSize: 'sm',
            color: 'faintText',
            letterSpacing: '0.005em',
            _selected: {
              color: 'mainText',
              fontWeight: 600,
            },
            _hover: { color: 'dimText' },
          },
          tablist: {
            borderColor: 'defaultBorder',
          },
          tabindicator: {
            bg: 'brand.600',
            height: '2px',
          },
        },
      },
    },
    Select: {
      variants: {
        outline: {
          field: {
            borderRadius: '10px',
            bg: 'inputBg',
            borderColor: 'defaultBorder',
          },
        },
      },
    },
    Menu: {
      baseStyle: {
        list: {
          borderRadius: 'xl',
          border: '1px solid',
          borderColor: 'defaultBorder',
          bg: 'cardBg',
          boxShadow: '0 4px 24px rgba(15,23,42,0.10)',
          overflow: 'hidden',
          py: 1,
          _dark: {
            boxShadow: '0 4px 24px rgba(0,0,0,0.40)',
          },
        },
        item: {
          fontSize: 'sm',
          fontWeight: 450,
          color: 'dimText',
          bg: 'cardBg',
          _hover: { bg: 'subtleBg', color: 'mainText' },
        },
      },
    },
  },
})

export default theme
