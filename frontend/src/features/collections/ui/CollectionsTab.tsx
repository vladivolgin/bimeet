import {
  VStack,
  HStack,
  Text,
  Button,
  Box,
  Grid,
  IconButton,
  Input,
  NumberInput,
  NumberInputField,
  FormControl,
  FormLabel,
  useDisclosure,
  Collapse,
  Avatar,
  Tooltip,
  Icon,
  Progress,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  ModalCloseButton,
  Badge,
  Link,
} from '@chakra-ui/react'
import { useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { FiTrash2, FiPlus, FiCheck, FiPocket, FiExternalLink, FiUpload } from 'react-icons/fi'
import type { Event } from '@/entities/event/model/types'
import type { CollectionContribution } from '@/entities/collection/model/types'
import {
  useAddCollection,
  useRemoveCollection,
  useSubmitContribution,
  useConfirmContribution,
  useRejectContribution,
  useMarkPaid,
} from '@/features/collections/model/hooks'
import { useAuth } from '@/features/auth/model/AuthContext'

function ContributionStatusBadge({ status }: { status: CollectionContribution['status'] }) {
  const { t } = useTranslation()
  if (status === 'paid') return <Badge colorScheme="green">{t('collections.statusPaid')}</Badge>
  if (status === 'pending') return <Badge colorScheme="yellow">{t('collections.statusPending')}</Badge>
  return null
}

interface ReceiptUploadModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (file: File) => void
  isLoading: boolean
}

function ReceiptUploadModal({ isOpen, onClose, onSubmit, isLoading }: ReceiptUploadModalProps) {
  const { t } = useTranslation()
  const [file, setFile] = useState<File | null>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  function handleClose() {
    setFile(null)
    onClose()
  }

  function handleSubmit() {
    if (file) {
      onSubmit(file)
    }
  }

  return (
    <Modal isOpen={isOpen} onClose={handleClose} isCentered>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>{t('collections.uploadReceiptTitle')}</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <VStack spacing={4}>
            <Text fontSize="sm" color="dimText">{t('collections.uploadReceiptHint')}</Text>
            <input
              ref={inputRef}
              type="file"
              accept="image/*,.pdf"
              style={{ display: 'none' }}
              onChange={e => setFile(e.target.files?.[0] ?? null)}
            />
            <Button
              w="full"
              variant="outline"
              leftIcon={<FiUpload />}
              onClick={() => inputRef.current?.click()}
            >
              {file ? file.name : t('collections.chooseFile')}
            </Button>
          </VStack>
        </ModalBody>
        <ModalFooter gap={2}>
          <Button variant="ghost" onClick={handleClose}>{t('common.cancel')}</Button>
          <Button
            colorScheme="blue"
            isDisabled={!file}
            isLoading={isLoading}
            onClick={handleSubmit}
          >
            {t('collections.uploadReceiptButton')}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default function CollectionsTab({ event }: { event: Event }) {
  const { user } = useAuth()
  const { t } = useTranslation()
  const { isOpen, onToggle } = useDisclosure()
  const [newTitle, setNewTitle] = useState('')
  const [newAmount, setNewAmount] = useState('')
  const [activeCollectionId, setActiveCollectionId] = useState<string | null>(null)

  const addCollection      = useAddCollection(event.id)
  const removeCollection   = useRemoveCollection(event.id)
  const submitContribution = useSubmitContribution(event.id)
  const confirmContribution = useConfirmContribution(event.id)
  const rejectContribution  = useRejectContribution(event.id)
  const markPaid            = useMarkPaid(event.id)

  const isOrganizer = event.organizer.id === user?.id
  const isCompleted = event.status === 'completed'

  const confirmedParticipants = event.participants.filter(p => p.status === 'confirmed')
  const isOrganizerParticipant = confirmedParticipants.some(p => p.user.id === event.organizer.id)
  const contributors = isOrganizerParticipant
    ? confirmedParticipants
    : [
        ...confirmedParticipants,
        { id: 'organizer', user: event.organizer, status: 'confirmed' as const },
      ]

  const totalTarget = event.collections.reduce((s, c) => s + c.per_person_amount * contributors.length, 0)
  const totalPaidCount = event.collections.flatMap(c => c.contributions).filter(c => c.status === 'paid').length
  const totalSlots = event.collections.length * contributors.length
  const collectedPct = totalSlots > 0 ? Math.round((totalPaidCount / totalSlots) * 100) : 0
  const perPersonTotal = event.collections.reduce((s, c) => s + c.per_person_amount, 0)

  function handleAdd() {
    if (!newTitle.trim() || !newAmount) return
    addCollection.mutate({ title: newTitle.trim(), per_person_amount: Number(newAmount) }, {
      onSuccess: () => {
        setNewTitle('')
        setNewAmount('')
        onToggle()
      },
    })
  }

  function handleUploadReceipt(file: File) {
    if (!activeCollectionId) return
    submitContribution.mutate(
      { collectionId: activeCollectionId, file },
      { onSuccess: () => setActiveCollectionId(null) },
    )
  }

  function avatarColor(status: CollectionContribution['status']) {
    if (status === 'paid') return 'green.400'
    if (status === 'pending') return 'yellow.400'
    return 'faintText'
  }

  function avatarOpacity(status: CollectionContribution['status']) {
    return status === 'not_paid' ? 0.35 : 1
  }

  function tooltipLabel(name: string, status: CollectionContribution['status']) {
    if (status === 'paid') return t('collections.tooltipPaid', { name })
    if (status === 'pending') return t('collections.tooltipPending', { name })
    return t('collections.tooltipNotPaid', { name })
  }

  return (
    <VStack align="stretch" spacing={5}>
      <ReceiptUploadModal
        isOpen={activeCollectionId !== null}
        onClose={() => setActiveCollectionId(null)}
        onSubmit={handleUploadReceipt}
        isLoading={submitContribution.isPending}
      />

      {event.collections.length > 0 && (
        <Grid templateColumns="repeat(3, 1fr)" gap={2}>
          <Box p={3.5} borderRadius="xl" bg="subtleBg" border="1px solid" borderColor="subtleBorder" borderLeft="3px solid" borderLeftColor="purple.400">
            <Text fontSize="xs" color="dimText" fontWeight="500" mb={0.5}>{t('collections.totalLabel')}</Text>
            <Text fontSize="lg" fontWeight="700" color="mainText" letterSpacing="-0.5px">{totalTarget.toLocaleString('ru-RU')} ₽</Text>
          </Box>
          <Box p={3.5} borderRadius="xl" bg="subtleBg" border="1px solid" borderColor="subtleBorder" borderLeft="3px solid" borderLeftColor="brand.400">
            <Text fontSize="xs" color="dimText" fontWeight="500" mb={0.5}>{t('collections.perPersonLabel')}</Text>
            <Text fontSize="lg" fontWeight="700" color="mainText" letterSpacing="-0.5px">{perPersonTotal.toLocaleString('ru-RU')} ₽</Text>
            <Text fontSize="xs" color="dimText">{t('collections.participantsCount', { count: confirmedParticipants.length })}</Text>
          </Box>
          <Box p={3.5} borderRadius="xl" bg="subtleBg" border="1px solid" borderColor="subtleBorder" borderLeft="3px solid" borderLeftColor="green.400">
            <Text fontSize="xs" color="dimText" fontWeight="500" mb={0.5}>{t('collections.contributedLabel')}</Text>
            <Text fontSize="lg" fontWeight="700" color="mainText" letterSpacing="-0.5px">{collectedPct}%</Text>
            <Text fontSize="xs" color="dimText">{totalPaidCount}/{totalSlots}</Text>
          </Box>
        </Grid>
      )}

      {totalSlots > 0 && (
        <Box>
          <HStack justify="space-between" mb={1.5}>
            <Text fontSize="xs" color="faintText" fontWeight="500">{t('collections.progressLabel')}</Text>
            <Text fontSize="xs" color="brand.600" fontWeight="600">{collectedPct}%</Text>
          </HStack>
          <Progress value={collectedPct} borderRadius="full" size="xs" colorScheme="purple" />
        </Box>
      )}

      <HStack justify="space-between">
        <Text fontSize="xs" fontWeight="600" color="faintText" textTransform="uppercase" letterSpacing="0.06em">
          {t('collections.sectionTitle')}
        </Text>
        {isOrganizer && !isCompleted && (
          <Button size="sm" leftIcon={<FiPlus />} colorScheme="blue" variant="outline" onClick={onToggle}>
            Добавить
          </Button>
        )}
      </HStack>

      <Collapse in={isOpen} animateOpacity>
        <Box p={4} borderRadius="xl" bg="subtleBg" border="1px solid" borderColor="subtleBorder">
          <VStack spacing={3}>
            <FormControl>
              <FormLabel fontSize="xs" fontWeight="600" color="dimText">{t('collections.fieldTitle')}</FormLabel>
              <Input size="sm" value={newTitle} onChange={e => setNewTitle(e.target.value)} placeholder={t('collections.fieldTitlePlaceholder')} autoFocus />
            </FormControl>
            <FormControl>
              <FormLabel fontSize="xs" fontWeight="600" color="dimText">{t('collections.fieldAmount')}</FormLabel>
              <NumberInput size="sm" value={newAmount} onChange={setNewAmount} min={1}>
                <NumberInputField placeholder="0" />
              </NumberInput>
              {newAmount && confirmedParticipants.length > 0 && (
                <Text fontSize="xs" color="dimText" mt={1}>
                  {t('collections.perPersonHint', { amount: (Number(newAmount) * confirmedParticipants.length).toLocaleString('ru-RU') })}
                </Text>
              )}
            </FormControl>
            <HStack w="full" justify="flex-end" spacing={2}>
              <Button size="sm" variant="ghost" onClick={onToggle}>{t('common.cancel')}</Button>
              <Button size="sm" colorScheme="blue" onClick={handleAdd} isLoading={addCollection.isPending}>{t('common.create')}</Button>
            </HStack>
          </VStack>
        </Box>
      </Collapse>

      {event.collections.length === 0 && (
        <Box textAlign="center" py={10}>
          <Icon as={FiPocket} boxSize={8} color="faintText" mb={2} />
          <Text fontSize="sm" color="dimText">{t('collections.empty')}</Text>
        </Box>
      )}

      <VStack align="stretch" spacing={3}>
        {event.collections.map(collection => {
          const contributions = collection.contributions ?? []
          const paidCount = contributions.filter(c => c.status === 'paid').length
          const pendingContribs = contributions.filter(c => c.status === 'pending')
          const myContribution = contributions.find(c => c.user_id === user?.id)
          const allPaid = paidCount === contributors.length && contributors.length > 0
          const shareAmount = collection.per_person_amount
          const pct = contributors.length > 0 ? Math.round((paidCount / contributors.length) * 100) : 0

          return (
            <Box
              key={collection.id}
              p={4}
              borderRadius="xl"
              border="1px solid"
              borderColor={allPaid ? 'green.300' : 'subtleBorder'}
              bg="cardBg"
              transition="all 0.2s"
            >
              <HStack justify="space-between" mb={1}>
                <Box minW={0} flex={1}>
                  <Text fontWeight="500" fontSize="sm" noOfLines={1}>{collection.title}</Text>
                  <HStack spacing={1.5} mt={0.5}>
                    <Text fontSize="xs" color="dimText">{t('collections.totalAmount', { amount: (collection.per_person_amount * contributors.length).toLocaleString('ru-RU') })}</Text>
                    <Text fontSize="xs" color="faintText">·</Text>
                    <Text fontSize="xs" color="brand.600" fontWeight="600">{t('collections.perYou', { amount: shareAmount.toLocaleString('ru-RU') })}</Text>
                  </HStack>
                </Box>
                <HStack spacing={2} flexShrink={0}>
                  <Box px={2.5} py={1} borderRadius="md"
                    bg={allPaid ? 'green.100' : 'orange.50'}
                    _dark={{ bg: allPaid ? 'rgba(74,222,128,0.15)' : 'rgba(251,146,60,0.12)' }}
                  >
                    <Text fontSize="xs" fontWeight="600" color={allPaid ? 'green.600' : 'orange.500'}>{paidCount}/{contributors.length}</Text>
                  </Box>
                  {isOrganizer && !isCompleted && (
                    <Tooltip label={paidCount > 0 ? t('collections.deleteTooltip') : ''} isDisabled={paidCount === 0} borderRadius="lg">
                      <IconButton
                        aria-label={t('collections.deleteButton')}
                        icon={<FiTrash2 />}
                        size="sm"
                        variant="ghost"
                        colorScheme="red"
                        isDisabled={paidCount > 0}
                        onClick={() => removeCollection.mutate(collection.id)}
                      />
                    </Tooltip>
                  )}
                </HStack>
              </HStack>

              <Progress value={pct} borderRadius="full" size="xs" colorScheme={allPaid ? 'green' : 'purple'} my={3} />

              {/* Avatar row */}
              {contributors.length > 0 && (
                <HStack spacing={1.5} flexWrap="wrap" mb={3}>
                  {contributors.map(p => {
                    const contrib = contributions.find(c => c.user_id === p.user.id)
                    const status = contrib?.status ?? 'not_paid'
                    return (
                      <Tooltip key={p.user.id} label={tooltipLabel(p.user.name, status)} borderRadius="lg">
                        <Box position="relative">
                          <Avatar
                            size="xs"
                            name={p.user.name}
                            src={p.user.avatar_url ?? undefined}
                            opacity={avatarOpacity(status)}
                            bg={avatarColor(status)}
                            transition="all 0.2s"
                          />
                          {status === 'paid' && (
                            <Box position="absolute" bottom="-1px" right="-1px" w="9px" h="9px" bg="green.400" borderRadius="full" border="1.5px solid" borderColor="cardBg" />
                          )}
                          {status === 'pending' && (
                            <Box position="absolute" bottom="-1px" right="-1px" w="9px" h="9px" bg="yellow.400" borderRadius="full" border="1.5px solid" borderColor="cardBg" />
                          )}
                        </Box>
                      </Tooltip>
                    )
                  })}
                </HStack>
              )}

              {/* Participant action button */}
              {user && contributors.find(p => p.user.id === user.id) && !isCompleted && (
                <Box>
                  {/* Organizer pays directly (no receipt) */}
                  {isOrganizer && (!myContribution || myContribution.status === 'not_paid') && (
                    <Button
                      size="sm" w="full" colorScheme="blue" variant="outline"
                      leftIcon={<FiCheck />}
                      isLoading={markPaid.isPending}
                      onClick={() => markPaid.mutate({ collectionId: collection.id, userId: user.id })}
                    >
                      {t('collections.buttonContribute')}
                    </Button>
                  )}
                  {/* Regular participant — upload receipt */}
                  {!isOrganizer && (!myContribution || myContribution.status === 'not_paid') && (
                    <Button size="sm" w="full" variant="outline" leftIcon={<FiUpload />} onClick={() => setActiveCollectionId(collection.id)}>
                      {t('collections.buttonContribute')}
                    </Button>
                  )}
                  {myContribution?.status === 'pending' && (
                    <HStack>
                      <Badge colorScheme="yellow" px={3} py={1} borderRadius="md">{t('collections.statusPending')}</Badge>
                    </HStack>
                  )}
                  {myContribution?.status === 'paid' && (
                    <HStack>
                      <Icon as={FiCheck} color="green.500" />
                      <Badge colorScheme="green" px={3} py={1} borderRadius="md">{t('collections.statusPaid')}</Badge>
                    </HStack>
                  )}
                </Box>
              )}

              {/* Organizer: pending section */}
              {isOrganizer && !isCompleted && pendingContribs.length > 0 && (
                <Box mt={3} p={3} borderRadius="lg" bg="subtleBg" border="1px solid" borderColor="yellow.300" _dark={{ borderColor: 'yellow.600' }}>
                  <Text fontSize="xs" fontWeight="600" color="dimText" mb={2}>{t('collections.pendingSection')} ({pendingContribs.length})</Text>
                  <VStack align="stretch" spacing={2}>
                    {pendingContribs.map(contrib => {
                      const participant = contributors.find(p => p.user.id === contrib.user_id)
                      return (
                        <HStack key={contrib.id} justify="space-between">
                          <HStack spacing={2}>
                            <Avatar size="xs" name={participant?.user.name ?? '?'} src={participant?.user.avatar_url ?? undefined} />
                            <Text fontSize="xs" fontWeight="500">{participant?.user.name ?? contrib.user_id}</Text>
                            {contrib.receipt_url && (
                              <Link href={contrib.receipt_url} isExternal fontSize="xs" color="brand.500">
                                <HStack spacing={1}>
                                  <Icon as={FiExternalLink} boxSize={3} />
                                  <Text>{t('collections.viewReceipt')}</Text>
                                </HStack>
                              </Link>
                            )}
                          </HStack>
                          <HStack spacing={1}>
                            <Button
                              size="xs"
                              colorScheme="green"
                              isLoading={confirmContribution.isPending}
                              onClick={() => confirmContribution.mutate({ collectionId: collection.id, contribId: contrib.id })}
                            >
                              {t('collections.confirm')}
                            </Button>
                            <Button
                              size="xs"
                              colorScheme="red"
                              variant="outline"
                              isLoading={rejectContribution.isPending}
                              onClick={() => rejectContribution.mutate({ collectionId: collection.id, contribId: contrib.id })}
                            >
                              {t('collections.reject')}
                            </Button>
                          </HStack>
                        </HStack>
                      )
                    })}
                  </VStack>
                </Box>
              )}

              {/* Organizer: mark-paid for other participants (not themselves) */}
              {isOrganizer && !isCompleted && (
                <VStack align="stretch" spacing={1} mt={3}>
                  {contributors
                    .filter(p => {
                      if (p.user.id === user?.id) return false // organizer uses own button above
                      const contrib = contributions.find(c => c.user_id === p.user.id)
                      const status = contrib?.status ?? 'not_paid'
                      return status === 'not_paid' // pending shown in the block above, paid already done
                    })
                    .map(p => {
                      const contrib = contributions.find(c => c.user_id === p.user.id)
                      return (
                        <HStack key={p.user.id} justify="space-between">
                          <HStack spacing={2}>
                            <Avatar size="xs" name={p.user.name} src={p.user.avatar_url ?? undefined} />
                            <Text fontSize="xs">{p.user.name}</Text>
                            <ContributionStatusBadge status={contrib?.status ?? 'not_paid'} />
                          </HStack>
                          <Button
                            size="xs"
                            variant="ghost"
                            colorScheme="green"
                            isLoading={markPaid.isPending}
                            onClick={() => markPaid.mutate({ collectionId: collection.id, userId: p.user.id })}
                          >
                            {t('collections.markPaid')}
                          </Button>
                        </HStack>
                      )
                    })}
                </VStack>
              )}
            </Box>
          )
        })}
      </VStack>
    </VStack>
  )
}
