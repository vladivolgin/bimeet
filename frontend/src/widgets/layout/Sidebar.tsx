import {
  Box,
  HStack,
  Text,
  VStack,
  Avatar,
  Icon,
  Divider,
  useColorModeValue,
} from '@chakra-ui/react'
import { Link, useMatch, useNavigate } from 'react-router-dom'
import { FiCalendar, FiLogOut, FiUser } from 'react-icons/fi'
import { useTranslation } from 'react-i18next'
import { useAuth } from '@/features/auth/model/AuthContext'

interface NavItemProps {
  to: string
  icon: React.ElementType
  label: string
  onClick?: () => void
}

function NavItem({ to, icon, label, onClick }: NavItemProps) {
  const isActive = !!useMatch({ path: to, end: false })

  return (
    <Box
      as={Link}
      to={to}
      onClick={onClick}
      display="flex"
      alignItems="center"
      gap={3}
      px={3}
      py={2.5}
      borderRadius="10px"
      bg={isActive ? 'navActiveBg' : 'transparent'}
      color={isActive ? 'navActiveText' : 'dimText'}
      fontWeight={isActive ? 600 : 450}
      fontSize="sm"
      transition="all 0.15s"
      _hover={isActive ? {} : { bg: 'subtleBg', color: 'mainText' }}
      userSelect="none"
    >
      <Icon as={icon} boxSize={4} />
      {label}
    </Box>
  )
}

interface SidebarContentProps {
  onClose?: () => void
}

export default function SidebarContent({ onClose }: SidebarContentProps) {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const { t } = useTranslation()
  const logoutHoverBg = useColorModeValue('red.50', 'rgba(254,202,202,0.08)')

  function handleLogout() {
    logout()
    navigate('/login')
    onClose?.()
  }

  return (
    <Box
      h="100%"
      display="flex"
      flexDirection="column"
      bg="sidebarBg"
      borderRight="1px solid"
      borderColor="sidebarBorder"
    >
      {/* Logo */}
      <Box px={4} pt={5} pb={4}>
        <HStack spacing={2}>
          <Box
            w="28px"
            h="28px"
            borderRadius="md"
            bg="brand.600"
            display="flex"
            alignItems="center"
            justifyContent="center"
            flexShrink={0}
          >
            <Text fontSize="12px" lineHeight={1} color="white">✦</Text>
          </Box>
          <Text fontSize="md" fontWeight="700" color="mainText" letterSpacing="-0.3px">
            Bimeet
          </Text>
        </HStack>
      </Box>

      <Divider borderColor="sidebarBorder" />

      {/* Navigation */}
      <VStack align="stretch" spacing={0.5} px={3} pt={3} flex={1}>
        <NavItem to="/events" icon={FiCalendar} label={t('nav.events')} onClick={onClose} />
        <NavItem to="/profile" icon={FiUser} label={t('nav.profile')} onClick={onClose} />
      </VStack>

      {/* User section */}
      {user && (
        <Box px={3} pb={4}>
          <Divider borderColor="sidebarBorder" mb={3} />
          <Box
            as={Link}
            to="/profile"
            onClick={onClose}
            display="block"
            px={3}
            py={2}
            mb={1}
            borderRadius="10px"
            _hover={{ bg: 'subtleBg' }}
            transition="background 0.15s"
          >
            <HStack spacing={2.5}>
              <Avatar size="sm" src={user.avatar_url ?? undefined} name={user.name} bg="brand.600" color="white" fontWeight="600" flexShrink={0} />
              <Box minW={0}>
                <Text fontSize="sm" fontWeight="600" color="mainText" noOfLines={1}>{user.name}</Text>
                <Text fontSize="xs" color="faintText" noOfLines={1}>{user.email}</Text>
              </Box>
            </HStack>
          </Box>
          <Box
            display="flex"
            alignItems="center"
            gap={3}
            px={3}
            py={2.5}
            borderRadius="10px"
            color="dimText"
            fontSize="sm"
            fontWeight={450}
            cursor="pointer"
            transition="all 0.15s"
            _hover={{ bg: logoutHoverBg, color: 'red.500' }}
            onClick={handleLogout}
          >
            <Icon as={FiLogOut} boxSize={4} />
            {t('common.logout', 'Выйти')}
          </Box>
        </Box>
      )}
    </Box>
  )
}
