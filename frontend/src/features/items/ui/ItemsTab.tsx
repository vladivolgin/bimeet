import {
  VStack,
  HStack,
  Text,
  Button,
  Box,
  Avatar,
  Badge,
  Input,
  FormControl,
  FormLabel,
  useDisclosure,
  Collapse,
  Icon,
} from '@chakra-ui/react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { FiPlus, FiShoppingBag } from 'react-icons/fi'
import type { Event } from '@/entities/event/model/types'
import { useAddItem, useAssignItem } from '@/features/items/model/hooks'
import { useAuth } from '@/features/auth/model/AuthContext'

export default function ItemsTab({ event }: { event: Event }) {
  const { user } = useAuth()
  const { t } = useTranslation()
  const { isOpen, onToggle } = useDisclosure()
  const [newName, setNewName] = useState('')

  const addItem = useAddItem(event.id)
  const assignItem = useAssignItem(event.id)
  const isCompleted = event.status === 'completed'

  function handleAdd() {
    if (!newName.trim()) return
    addItem.mutate(newName.trim(), {
      onSuccess: () => {
        setNewName('')
        onToggle()
      },
    })
  }

  return (
    <VStack align="stretch" spacing={4}>
      <HStack justify="space-between">
        <Text fontSize="xs" fontWeight="600" color="faintText" textTransform="uppercase" letterSpacing="0.06em">
          {t('items.sectionTitle')}
        </Text>
        {!isCompleted && (
          <Button size="sm" leftIcon={<FiPlus />} colorScheme="blue" variant="outline" onClick={onToggle}>
            Добавить
          </Button>
        )}
      </HStack>

      <Collapse in={isOpen} animateOpacity>
        <Box p={4} borderRadius="xl" bg="subtleBg" border="1px solid" borderColor="subtleBorder">
          <VStack spacing={3}>
            <FormControl>
              <FormLabel fontSize="xs" fontWeight="600" color="dimText">{t('items.fieldName')}</FormLabel>
              <Input
                size="sm"
                value={newName}
                onChange={e => setNewName(e.target.value)}
                placeholder={t('items.fieldNamePlaceholder')}
                onKeyDown={e => e.key === 'Enter' && handleAdd()}
                autoFocus
              />
            </FormControl>
            <HStack w="full" justify="flex-end" spacing={2}>
              <Button size="sm" variant="ghost" onClick={onToggle}>{t('common.cancel')}</Button>
              <Button size="sm" colorScheme="blue" onClick={handleAdd} isLoading={addItem.isPending}>{t('common.add')}</Button>
            </HStack>
          </VStack>
        </Box>
      </Collapse>

      {event.items.length === 0 && (
        <Box textAlign="center" py={10}>
          <Icon as={FiShoppingBag} boxSize={8} color="faintText" mb={2} />
          <Text fontSize="sm" color="dimText">{t('items.empty')}</Text>
        </Box>
      )}

      <VStack align="stretch" spacing={2}>
        {event.items.map(item => {
          const isMine = item.assigned_to?.id === user?.id
          const isTaken = Boolean(item.assigned_to)

          return (
            <Box
              key={item.id}
              p={3.5}
              borderRadius="xl"
              border="1px solid"
              borderColor="subtleBorder"
              bg="cardBg"
            >
              <HStack justify="space-between" mb={1}>
                <HStack spacing={3} minW={0} flex={1}>
                  {item.assigned_to ? (
                    <Avatar size="sm" name={item.assigned_to.name} bg="brand.400" color="white" flexShrink={0} />
                  ) : (
                    <Box w="32px" h="32px" borderRadius="full" bg="subtleBg" border="1px solid" borderColor="subtleBorder" flexShrink={0} />
                  )}
                  <Box minW={0}>
                    <Text fontSize="sm" fontWeight="500" noOfLines={1}>{item.name}</Text>
                    <Text fontSize="xs" color="dimText" mt={0.5}>
                      {item.assigned_to ? item.assigned_to.name : t('items.unassigned')}
                    </Text>
                  </Box>
                </HStack>

                <HStack spacing={2} flexShrink={0}>
                  {isMine && (
                    <Badge colorScheme="blue" variant="subtle">{t('items.badgeYou')}</Badge>
                  )}
                  {isCompleted ? (
                    isTaken && !isMine && <Badge colorScheme="green" variant="subtle">{t('items.badgeTaken')}</Badge>
                  ) : isMine ? (
                    <Button
                      size="sm"
                      colorScheme="red"
                      variant="ghost"
                      onClick={() => assignItem.mutate({ itemId: item.id, userId: null })}
                      isLoading={assignItem.isPending}
                    >
                      {t('items.buttonGiveUp')}
                    </Button>
                  ) : !isTaken ? (
                    <Button
                      size="sm"
                      colorScheme="blue"
                      variant="outline"
                      onClick={() => user && assignItem.mutate({ itemId: item.id, userId: user.id })}
                      isLoading={assignItem.isPending}
                    >
                      {t('items.buttonTake')}
                    </Button>
                  ) : (
                    <Badge colorScheme="green" variant="subtle">{t('items.badgeTaken')}</Badge>
                  )}
                </HStack>
              </HStack>
            </Box>
          )
        })}
      </VStack>
    </VStack>
  )
}
