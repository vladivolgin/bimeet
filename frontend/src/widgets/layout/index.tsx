import {
  Box,
  Button,
  Drawer,
  DrawerBody,
  DrawerContent,
  DrawerOverlay,
  IconButton,
  Tooltip,
  useColorMode,
  useColorModeValue,
  useDisclosure,
} from '@chakra-ui/react'
import { FiMenu, FiSun, FiMoon } from 'react-icons/fi'
import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'
import SidebarContent from './Sidebar'
import BottomNav from './BottomNav'
import NotificationCenter from '@/features/notifications/ui/NotificationCenter'

function LanguageToggle() {
  const { i18n } = useTranslation()
  const isRu = i18n.language === 'ru'
  const hoverBg = useColorModeValue('rgba(15,23,42,0.05)', 'rgba(255,255,255,0.07)')
  return (
    <Button
      size="xs"
      variant="ghost"
      fontWeight="600"
      fontSize="xs"
      color="dimText"
      borderRadius="8px"
      _hover={{ color: 'mainText', bg: hoverBg }}
      onClick={() => i18n.changeLanguage(isRu ? 'en' : 'ru')}
    >
      {isRu ? 'EN' : 'RU'}
    </Button>
  )
}

export default function Layout({ children }: { children: ReactNode }) {
  const { colorMode, toggleColorMode } = useColorMode()
  const { t } = useTranslation()
  const { isOpen, onOpen, onClose } = useDisclosure()
  const topBarBg = useColorModeValue('rgba(255,255,255,0.85)', 'rgba(19,20,31,0.90)')
  const toggleHoverBg = useColorModeValue('rgba(15,23,42,0.05)', 'rgba(255,255,255,0.07)')

  return (
    <Box minH="100vh" bg="pageBg" display="flex">
      {/* Sidebar — desktop only */}
      <Box
        display={{ base: 'none', lg: 'block' }}
        w="240px"
        flexShrink={0}
        position="sticky"
        top={0}
        h="100vh"
        overflowY="auto"
      >
        <SidebarContent />
      </Box>

      {/* Mobile drawer */}
      <Drawer isOpen={isOpen} placement="left" onClose={onClose} size="xs">
        <DrawerOverlay />
        <DrawerContent bg="sidebarBg" maxW="240px" p={0}>
          <DrawerBody p={0}>
            <SidebarContent onClose={onClose} />
          </DrawerBody>
        </DrawerContent>
      </Drawer>

      {/* Content area */}
      <Box flex={1} minW={0} display="flex" flexDirection="column">
        {/* TopBar */}
        <Box
          position="sticky"
          top={0}
          zIndex={50}
          backdropFilter="blur(16px) saturate(1.6)"
          bg={topBarBg}
          borderBottom="1px solid"
          borderColor="sidebarBorder"
          h="48px"
          display="flex"
          alignItems="center"
          justifyContent="flex-end"
          px={4}
        >
          {/* Hamburger — mobile only */}
          <IconButton
            aria-label={t('layout.openMenu')}
            icon={<FiMenu size={18} />}
            variant="ghost"
            size="sm"
            borderRadius="10px"
            onClick={onOpen}
            color="dimText"
            _hover={{ color: 'mainText', bg: toggleHoverBg }}
            display={{ base: 'flex', lg: 'none' }}
            mr="auto"
          />

          <NotificationCenter />

          <LanguageToggle />

          <Tooltip label={colorMode === 'light' ? t('layout.darkTheme') : t('layout.lightTheme')} borderRadius="lg">
            <IconButton
              aria-label={t('layout.toggleTheme')}
              icon={colorMode === 'light' ? <FiMoon size={15} /> : <FiSun size={15} />}
              variant="ghost"
              size="sm"
              borderRadius="10px"
              onClick={toggleColorMode}
              color="dimText"
              _hover={{ color: 'mainText', bg: toggleHoverBg }}
            />
          </Tooltip>
        </Box>

        {/* Main content */}
        <Box
          as="main"
          flex={1}
          p={{ base: 4, sm: 6, xl: 8 }}
          pb={{ base: '72px', lg: 8 }}  // extra bottom padding for BottomNav on mobile
        >
          {children}
        </Box>
      </Box>

      {/* Bottom nav — mobile only */}
      <BottomNav />
    </Box>
  )
}
