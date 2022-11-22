# https://github.com/DMOJ/online-judge/blob/master/judge/bridge/judge_handler.py#L223
def submit(self, id, problem, language, source):
  data = self.get_related_submission_data(id)
  self._working = id
  self._no_response_job = threading.Timer(20, self._kill_if_no_response)
  self.send({
      'name': 'submission-request',
      'submission-id': id,
      'problem-id': problem,
      'language': language,
      'source': source,
      'time-limit': data.time,
      'memory-limit': data.memory,
      'short-circuit': data.short_circuit,
      'meta': {
          'pretests-only': data.pretests_only,
          'in-contest': data.contest_no,
          'attempt-no': data.attempt_no,
          'user': data.user_id,
      },
  })

# https://github.com/DMOJ/online-judge/blob/master/judge/bridge/judge_list.py#L118
def judge(self, id, problem, language, source, judge_id, priority):
    with self.lock:
        if id in self.submission_map or id in self.node_map:
            # Already judging, don't queue again. This can happen during batch rejudges, rejudges should be
            # idempotent.
            return

        candidates = [judge for judge in self.judges if judge.can_judge(problem, language, judge_id)]
        available = [judge for judge in candidates if not judge.working]
        if judge_id:
            logger.info('Specified judge %s is%savailable', judge_id, ' ' if available else ' not ')
        else:
            logger.info('Free judges: %d', len(available))

        if len(candidates) > 1 and len(available) == 1 and priority >= REJUDGE_PRIORITY:
            available = []

        if available:
            # Schedule the submission on the judge reporting least load.
            judge = min(available, key=lambda judge: (judge.load, random()))
            logger.info('Dispatched submission %d to: %s', id, judge.name)
            self.submission_map[id] = judge
            try:
                judge.submit(id, problem, language, source)
            except Exception:
                logger.exception('Failed to dispatch %d (%s, %s) to %s', id, problem, language, judge.name)
                self.judges.discard(judge)
                return self.judge(id, problem, language, source, judge_id, priority)
        else:
            self.node_map[id] = self.queue.insert(
                (id, problem, language, source, judge_id),
                self.priority[priority],
            )
            logger.info('Queued submission: %d', id)

# https://github.com/DMOJ/online-judge/blob/master/judge/bridge/django_handler.py#L35
def on_submission(self, data):
    id = data['submission-id']
    problem = data['problem-id']
    language = data['language']
    source = data['source']
    judge_id = data['judge-id']
    priority = data['priority']
    if not self.judges.check_priority(priority):
        return {'name': 'bad-request'}
    self.judges.judge(id, problem, language, source, judge_id, priority)
    return {'name': 'submission-received', 'submission-id': id}