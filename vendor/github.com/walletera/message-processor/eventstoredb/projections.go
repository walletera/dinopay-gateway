package eventstoredb

import (
    "context"
    "fmt"
    "net/http"
    "net/url"
    "strings"
)

const (
    eventStoreByCategorySeparator = "last\n."
)

func EnableByCategoryProjection(ctx context.Context, esdbUrl string) error {
    parsedUrl, err := url.Parse(esdbUrl)
    if err != nil {
        return err
    }

    req, err := http.NewRequestWithContext(
        ctx,
        http.MethodPost,
        fmt.Sprintf("http://%s:%s/projection/$by_category/command/enable", parsedUrl.Hostname(), parsedUrl.Port()),
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed creating request for enabling $by_category projection: %w", err)
    }

    req.Header.Add("Accept", "application/json")
    req.Header.Add("Content-Length", "0")

    _, err = http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed enabling $by_category projection: %w", err)
    }

    return nil
}

func SetESDBByCategoryProjectionSeparator(ctx context.Context, esdbUrl string) error {
    parsedUrl, err := url.Parse(esdbUrl)
    if err != nil {
        return err
    }

    req, err := http.NewRequestWithContext(
        ctx,
        http.MethodPut,
        fmt.Sprintf("http://%s:%s/projection/$by_category/query?emit=yes", parsedUrl.Hostname(), parsedUrl.Port()),
        strings.NewReader(eventStoreByCategorySeparator),
    )
    if err != nil {
        return fmt.Errorf("failed creating request to update byCategory projection separator: %w", err)
    }

    req.Header.Set("Content-Type", "application/json; charset=utf-8")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed creating request to update byCategory projection separator: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed creating request to update byCategory projection separator - response status code %d", resp.StatusCode)
    }

    return nil
}
