# Crony

Lightweight cron replacement that calls URL

## Usage

Config

```json
{
  "tasks": [
    {
      "cron": "@every 1m",
      "url": "http://app/task"
    }
  ]
}

```  

Docker

 
```yaml
services:
  app:
    
  cron:
    image: 'formapro/crony:latest'
    depends_on:
      - app
    volumes:
      - './crony.json:/app/crony.json'
```

Symfony

Route:

```yaml
#config/routes.yaml

cron_task:
  path: '/cron/task/foo'
  controller: 'App\Controller\CronTaskController::foo'
  methods: [POST]
```

Controller:

```php
<?php declare(strict_types=1);

namespace App\Controller;

use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Event\PostResponseEvent;
use Symfony\Component\HttpKernel\KernelEvents;

class CronTaskController implements EventSubscriberInterface
{
    public function foo(Request $request): Response
    {
        $request->attributes->set('task', 'foo');

        return new Response();
    }

    public function onTerminate(PostResponseEvent $event): void
    {
        if ('foo' !== $event->getRequest()->attributes->get('task')) {
            return;
        }

        // do task
    }

    public static function getSubscribedEvents(): array
    {
        return [
            KernelEvents::TERMINATE => 'onTerminate',
        ];
    }
}
```


## License

MIT
