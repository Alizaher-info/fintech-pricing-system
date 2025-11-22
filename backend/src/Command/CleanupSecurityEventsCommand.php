<?php

namespace App\Command;

use App\Repository\SecurityEventRepository;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Style\SymfonyStyle;

#[AsCommand(
    name: 'app:cleanup-security-events',
    description: 'Delete old security events based on retention policy',
)]
class CleanupSecurityEventsCommand extends Command
{
    private const DEFAULT_RETENTION_DAYS = 365; // 1 year

    public function __construct(
        private EntityManagerInterface $entityManager,
        private SecurityEventRepository $securityEventRepository
    ) {
        parent::__construct();
    }

    protected function configure(): void
    {
        $this
            ->addOption(
                'days',
                'd',
                InputOption::VALUE_REQUIRED,
                'Number of days to retain events',
                self::DEFAULT_RETENTION_DAYS
            )
            ->addOption(
                'dry-run',
                null,
                InputOption::VALUE_NONE,
                'Show what would be deleted without actually deleting'
            );
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $io = new SymfonyStyle($input, $output);
        $retentionDays = (int) $input->getOption('days');
        $dryRun = $input->getOption('dry-run');

        $io->title('Security Events Cleanup Task');
        $io->info(sprintf('Retention policy: %d days', $retentionDays));

        // Calculate cutoff date
        $cutoffDate = new \DateTimeImmutable(sprintf('-%d days', $retentionDays));
        $io->text(sprintf('Deleting events older than: %s', $cutoffDate->format('Y-m-d H:i:s')));

        try {
            // Count records to be deleted
            $qb = $this->entityManager->createQueryBuilder();
            $count = $qb->select('COUNT(se.id)')
                ->from('App\Entity\SecurityEvent', 'se')
                ->where('se.createdAt < :cutoffDate')
                ->setParameter('cutoffDate', $cutoffDate)
                ->getQuery()
                ->getSingleScalarResult();

            if ($count === 0) {
                $io->success('No security events to clean up.');
                return Command::SUCCESS;
            }

            $io->warning(sprintf('Found %d security events to delete', $count));

            if ($dryRun) {
                $io->note('DRY RUN MODE - No records will be deleted');
                
                // Show sample of events that would be deleted
                $samples = $this->entityManager->createQueryBuilder()
                    ->select('se.id', 'se.eventType', 'se.createdAt', 'se.severity')
                    ->from('App\Entity\SecurityEvent', 'se')
                    ->where('se.createdAt < :cutoffDate')
                    ->setParameter('cutoffDate', $cutoffDate)
                    ->setMaxResults(5)
                    ->getQuery()
                    ->getArrayResult();

                $io->table(
                    ['ID', 'Event Type', 'Created At', 'Severity'],
                    array_map(fn($s) => [
                        $s['id'],
                        $s['eventType'],
                        $s['createdAt']->format('Y-m-d H:i:s'),
                        $s['severity']
                    ], $samples)
                );

                $io->success(sprintf('DRY RUN: Would delete %d records', $count));
                return Command::SUCCESS;
            }

            // Perform deletion
            $deleted = $this->entityManager->createQueryBuilder()
                ->delete('App\Entity\SecurityEvent', 'se')
                ->where('se.createdAt < :cutoffDate')
                ->setParameter('cutoffDate', $cutoffDate)
                ->getQuery()
                ->execute();

            $io->success(sprintf('Successfully deleted %d security events', $deleted));

            // Log summary
            $io->section('Cleanup Summary');
            $io->definitionList(
                ['Retention Period' => sprintf('%d days', $retentionDays)],
                ['Cutoff Date' => $cutoffDate->format('Y-m-d H:i:s')],
                ['Records Deleted' => $deleted],
                ['Execution Time' => date('Y-m-d H:i:s')]
            );

            return Command::SUCCESS;

        } catch (\Exception $e) {
            $io->error('Failed to cleanup security events: ' . $e->getMessage());
            return Command::FAILURE;
        }
    }
}
