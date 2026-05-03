import { Box, HStack, VStack, Icon, Text } from '@chakra-ui/react'
import { Link, useMatch } from 'react-router-dom'
import { FiCalendar, FiUser } from 'react-icons/fi'
import { useTranslation } from 'react-i18next'

interface BottomNavItemProps {
  to: string
  icon: React.ElementType
  label: string
}

function BottomNavItem({ to, icon, label }: BottomNavItemProps) {
  const isActive = !!useMatch({ path: to, end: false })

  return (
    <Box
      as={Link}
      to={to}
      flex={1}
      display="flex"
      alignItems="center"
      justifyContent="center"
      py={2}
    >
      <VStack spacing={0.5}>
        <Icon
          as={icon}
          boxSize={5}
          color={isActive ? 'brand.600' : 'faintText'}
          transition="color 0.15s"
        />
        <Text
          fontSize="10px"
          fontWeight={isActive ? 600 : 400}
          color={isActive ? 'brand.600' : 'faintText'}
          transition="all 0.15s"
        >
          {label}
        </Text>
      </VStack>
    </Box>
  )
}

export default function BottomNav() {
  const { t } = useTranslation()
  return (
    <Box
      display={{ base: 'flex', lg: 'none' }}
      position="fixed"
      bottom={0}
      left={0}
      right={0}
      zIndex={100}
      bg="sidebarBg"
      borderTop="1px solid"
      borderColor="sidebarBorder"
      h="56px"
      alignItems="stretch"
    >
      <HStack w="full" spacing={0} align="stretch">
        <BottomNavItem to="/events" icon={FiCalendar} label={t('nav.events')} />
        <BottomNavItem to="/profile" icon={FiUser} label={t('nav.profile')} />
      </HStack>
    </Box>
  )
}
